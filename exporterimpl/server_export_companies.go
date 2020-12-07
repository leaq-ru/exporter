package exporterimpl

import (
	"context"
	"errors"
	"github.com/nnqq/scr-exporter/csv"
	"github.com/nnqq/scr-exporter/safeerr"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"github.com/nnqq/scr-proto/codegen/go/opts"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"golang.org/x/sync/errgroup"
	"io"
	"time"
)

func (s *server) ExportCompanies(
	ctx context.Context,
	req *parser.GetListRequest,
) (
	res *exporter.ExportCompaniesResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// TODO redis cache

	compStream, err := s.companyClient.GetFull(ctx, &parser.GetV2Request{
		Opts: &opts.Page{
			Limit: 1000,
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
	})
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	var eg errgroup.Group
	compCh := make(chan *parser.FullCompanyV2)
	eg.Go(func() (e error) {
		defer close(compCh)

		for {
			comp, er := compStream.Recv()
			e = er
			if e != nil {
				if errors.Is(e, io.EOF) {
					e = nil
				} else {
					s.logger.Error().Err(e).Send()
				}
				return
			}

			compCh <- comp
		}
	})

	var csvPath string
	eg.Go(func() (e error) {
		csvPath, e = csv.Create(compCh)
		return
	})
	err = eg.Wait()
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	s3URL, err := s.store.Put(ctx, csvPath, true)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	res = &exporter.ExportCompaniesResponse{
		Url: s3URL,
	}
	return
}
