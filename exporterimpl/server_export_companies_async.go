package exporterimpl

import (
	"context"
	"errors"
	"github.com/leaq-ru/exporter/file"
	"github.com/leaq-ru/exporter/md"
	"github.com/leaq-ru/exporter/safeerr"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
	"sort"
	"strings"
	"time"
)

func (s *server) ExportCompaniesAsync(
	ctx context.Context,
	req *parser.GetListRequest,
) (
	res *emptypb.Empty,
	err error,
) {
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

	var eg errgroup.Group
	const limitPrettyName = 2
	var cityNames []string
	if len(req.GetCityIds()) != 0 {
		eg.Go(func() error {
			resCity, e := s.cityClient.GetCityByIds(ctx, &parser.GetCityByIdsRequest{
				CityIds: cut(req.GetCityIds(), limitPrettyName),
			})
			if e != nil {
				return e
			}

			for _, item := range resCity.GetCities() {
				cityNames = append(cityNames, item.GetTitle())
			}
			return nil
		})
	}

	var categoryNames []string
	if len(req.GetCategoryIds()) != 0 {
		eg.Go(func() error {
			resCategory, e := s.categoryClient.GetCategoryByIds(ctx, &parser.GetCategoryByIdsRequest{
				CategoryIds: cut(req.GetCategoryIds(), limitPrettyName),
			})
			if e != nil {
				return e
			}

			for _, item := range resCategory.GetCategories() {
				categoryNames = append(categoryNames, item.GetTitle())
			}
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	var rawName []string

	if len(cityNames) != 0 {
		sort.Strings(cityNames)
	}
	for _, name := range cityNames {
		rawName = append(rawName, name)
	}

	if len(categoryNames) != 0 {
		sort.Strings(categoryNames)
	}
	for _, name := range categoryNames {
		rawName = append(rawName, name)
	}

	if len(req.GetCityIds())+len(req.GetCategoryIds()) > len(rawName) {
		rawName = append(rawName, "...")
	}

	name := "Выгрузка CSV"
	if len(rawName) != 0 {
		name = strings.Join(rawName, ", ")
	}

	eventOID := primitive.NewObjectID()

	err = s.fileModel.EnsureLimitAndCreatePending(ctx, authUserOID, eventOID, name)
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

	res = &emptypb.Empty{}
	return
}

func cut(in []string, maxLen int) (out []string) {
	if len(in) > maxLen {
		out = in[:maxLen]
	} else {
		out = in
	}
	return
}
