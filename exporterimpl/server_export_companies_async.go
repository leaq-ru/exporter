package exporterimpl

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/md"
	"github.com/nnqq/scr-exporter/safeerr"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s *server) ExportCompaniesAsync(ctx context.Context, req *parser.GetListRequest) (res *empty.Empty, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	premium, err := md.GetDataPremium(ctx)
	if err != nil {
		return
	}

	if !premium {
		err = errors.New("this method is allowed only on data-premium plan")
		return
	}

	authUserOID, err := md.GetUserOID(ctx)
	if err != nil {
		return
	}

	eventOID := primitive.NewObjectID()

	err = s.fileModel.EnsureLimitAndCreatePending(ctx, authUserOID, eventOID, "Выгрузка csv")
	if err != nil {
		if errors.Is(err, file.ErrConcExports) {
			return
		}

		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	err = s.processAsync(eventOID, req)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	res = &empty.Empty{}
	return
}
