package exporterimpl

import (
	"context"
	"errors"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/csv"
	"github.com/nnqq/scr-exporter/safeerr"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"github.com/nnqq/scr-proto/codegen/go/opts"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"time"
)

func (s *server) ExportCompanies(
	ctx context.Context,
	req *parser.GetListRequest,
) (
	res *exporter.ExportCompaniesResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	reqComp := &parser.GetV2Request{
		Opts: &opts.Page{
			Limit: 2500,
		},
		CityIds:            req.GetCityIds(),
		CategoryIds:        req.GetCategoryIds(),
		HasEmail:           req.GetHasEmail(),
		HasPhone:           req.GetHasPhone(),
		HasOnline:          req.GetHasOnline(),
		HasInn:             req.GetHasInn(),
		HasKpp:             req.GetHasKpp(),
		HasOgrn:            req.GetHasOgrn(),
		HasAppStore:        req.GetHasAppStore(),
		HasGooglePlay:      req.GetHasGooglePlay(),
		HasVk:              req.GetHasVk(),
		VkMembersCount:     req.GetVkMembersCount(),
		HasInstagram:       req.GetHasInstagram(),
		HasTwitter:         req.GetHasTwitter(),
		HasYoutube:         req.GetHasYoutube(),
		HasFacebook:        req.GetHasFacebook(),
		TechnologyIds:      req.GetTechnologyIds(),
		TechnologyFindRule: req.GetTechnologyFindRule(),
	}

	s3URL, err := s.cachedExportModel.Get(ctx, reqComp)
	if err != nil {
		if errors.Is(err, cached_export.ErrNoFound) {
			err = nil

			compStream, e := s.companyClient.GetFull(ctx, &parser.GetFullRequest{
				Query: reqComp,
			})
			if e != nil {
				s.logger.Error().Err(e).Send()
				err = safeerr.InternalServerError
				return
			}

			csvPath, e := csv.DoPipelineSync(ctx, compStream)
			if e != nil {
				s.logger.Error().Err(e).Send()
				err = safeerr.InternalServerError
				return
			}

			s3URL, e = s.exporterBucket.Put(ctx, csvPath, true)
			if e != nil {
				s.logger.Error().Err(e).Send()
				err = safeerr.InternalServerError
				return
			}

			e = s.cachedExportModel.Set(ctx, reqComp, s3URL)
			if e != nil {
				s.logger.Error().Err(e).Send()
				err = safeerr.InternalServerError
				return
			}
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
			return
		}
	}

	res = &exporter.ExportCompaniesResponse{
		Url: s3URL,
	}
	return
}
