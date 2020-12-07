package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/event_log"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-exporter/store"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	store store.Store,
	companyClient parser.CompanyClient,
	fileModel file.Model,
	eventLogModel event_log.Model,
	mongoStartSession mongo.StartSession,
	serviceName string,
) Consumer {
	return Consumer{
		logger:            logger,
		stanConn:          stanConn,
		store:             store,
		companyClient:     companyClient,
		fileModel:         fileModel,
		eventLogModel:     eventLogModel,
		mongoStartSession: mongoStartSession,
		serviceName:       serviceName,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
