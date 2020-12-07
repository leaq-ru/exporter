package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/csv"
	"github.com/nnqq/scr-proto/codegen/go/opts"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
	"io"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		// TODO redis cache

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
		defer cancel()

		ack := func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		var msg message
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("got malformed msg, just ack")
			ack()
			return
		}

		alreadyProcessed, err := c.eventLogModel.AlreadyProcessed(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		if alreadyProcessed {
			ack()
			return
		}

		if rawMsg.RedeliveryCount >= 15 {
			c.logger.Error().Msg("got dead letter message, set status=fail and ack")

			err = c.fileModel.SetFail(ctx, msg.ID)
			if err != nil {
				c.logger.Error().Err(err).Send()
				return
			}

			ack()
			return
		}

		err = c.fileModel.SetInProgress(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		compStream, err := c.companyClient.GetFull(ctx, &parser.GetV2Request{
			Opts: &opts.Page{
				Limit: 100000,
			},
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
		})

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
			c.logger.Error().Err(err).Send()
			return
		}

		s3URL, err := c.store.Put(ctx, csvPath, true)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		sess, err := c.mongoStartSession()
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		_, err = sess.WithTransaction(ctx, func(sc mongo.SessionContext) (_ interface{}, e error) {
			e = c.eventLogModel.Put(sc, msg.ID)
			if e != nil {
				return
			}

			e = c.fileModel.SetSuccess(sc, msg.ID, s3URL)
			return
		})
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		ack()
	}()
}
