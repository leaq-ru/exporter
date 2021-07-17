package exporterimpl

import (
	"context"
	"github.com/leaq-ru/exporter/md"
	"github.com/leaq-ru/exporter/pagination"
	"github.com/leaq-ru/exporter/safeerr"
	"github.com/leaq-ru/proto/codegen/go/exporter"
	"time"
)

func (s *server) GetMyFiles(
	ctx context.Context,
	req *exporter.GetMyFilesRequest,
) (
	res *exporter.GetMyFilesResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	authUserOID, err := md.GetUserOID(ctx)
	if err != nil {
		return
	}

	limit, err := pagination.ApplyDefaultLimit(req)
	if err != nil {
		return
	}

	files, err := s.fileModel.Get(ctx, authUserOID, req.GetOpts().GetSkip(), limit)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	res = &exporter.GetMyFilesResponse{}
	for _, f := range files {
		res.Files = append(res.Files, &exporter.File{
			Id:           f.ID.Hex(),
			Name:         f.Name,
			Url:          f.URL,
			CreatedAt:    f.CreatedAt.String(),
			Status:       exporter.Status(f.Status),
			CurrentCount: f.CurrentCount,
			TotalCount:   f.TotalCount,
		})
	}
	return
}
