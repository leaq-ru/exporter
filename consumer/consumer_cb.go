package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/csv"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"golang.org/x/sync/errgroup"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
		defer cancel()

		defer func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}()

		var msg message
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("got malformed msg, just ack")
			return
		}

		alreadyProcessed, err := c.eventLogModel.AlreadyProcessed(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		if alreadyProcessed {
			return
		}

		failed := func() {
			e := c.fileModel.SetFail(ctx, msg.ID)
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		err = c.eventLogModel.Put(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			failed()
			return
		}

		reqComp := &parser.GetV2Request{
			CityIds:            msg.Query.GetCityIds(),
			CategoryIds:        msg.Query.GetCategoryIds(),
			HasEmail:           msg.Query.GetHasEmail(),
			HasPhone:           msg.Query.GetHasPhone(),
			HasOnline:          msg.Query.GetHasOnline(),
			HasInn:             msg.Query.GetHasInn(),
			HasKpp:             msg.Query.GetHasKpp(),
			HasOgrn:            msg.Query.GetHasOgrn(),
			HasAppStore:        msg.Query.GetHasAppStore(),
			HasGooglePlay:      msg.Query.GetHasGooglePlay(),
			HasVk:              msg.Query.GetHasVk(),
			VkMembersCount:     msg.Query.GetVkMembersCount(),
			HasInstagram:       msg.Query.GetHasInstagram(),
			HasTwitter:         msg.Query.GetHasTwitter(),
			HasYoutube:         msg.Query.GetHasYoutube(),
			HasFacebook:        msg.Query.GetHasFacebook(),
			TechnologyIds:      msg.Query.GetTechnologyIds(),
			TechnologyFindRule: msg.Query.GetTechnologyFindRule(),
		}

		s3URL, err := c.cachedExportModel.Get(ctx, reqComp)
		if err != nil {
			if errors.Is(err, cached_export.ErrNoFound) {
				err = nil

				var eg errgroup.Group
				var totalCount uint32
				eg.Go(func() (egErr error) {
					resCount, errCount := c.companyClient.GetCount(ctx, reqComp)
					if errCount != nil {
						egErr = errCount
						return
					}

					totalCount = resCount.GetCount()
					return
				})

				var compStream parser.Company_GetFullClient
				eg.Go(func() (egErr error) {
					compStream, egErr = c.companyClient.GetFull(ctx, reqComp)
					return
				})
				e := eg.Wait()
				if e != nil {
					c.logger.Error().Err(e).Send()
					failed()
					return
				}

				e = c.fileModel.SetInProgress(ctx, msg.ID, totalCount)
				if e != nil {
					c.logger.Error().Err(e).Send()
					failed()
					return
				}

				var currentCount uint32
				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							time.Sleep(10 * time.Second)
							errSetCount := c.fileModel.SetCurrentCount(ctx, msg.ID, currentCount)
							if errSetCount != nil {
								c.logger.Error().Err(errSetCount).Send()
							}
						}
					}
				}()

				csvPath, e := csv.DoPipeline(ctx, compStream, &currentCount)
				if e != nil {
					c.logger.Error().Err(e).Send()
					failed()
					return
				}

				s3URL, e = c.exporterBucket.Put(ctx, csvPath, true)
				if e != nil {
					c.logger.Error().Err(e).Send()
					failed()
					return
				}

				e = c.cachedExportModel.Set(ctx, reqComp, s3URL)
				if e != nil {
					c.logger.Error().Err(e).Send()
					failed()
					return
				}
			} else {
				c.logger.Error().Err(err).Send()
				failed()
				return
			}
		}

		err = c.fileModel.SetSuccess(ctx, msg.ID, s3URL)
		if err != nil {
			c.logger.Error().Err(err).Send()
			failed()
		}
		return
	}()
}
