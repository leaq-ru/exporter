package exporter_async

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/event_log"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stan stan.Conn,
	fileModel file.Model,
	eventLogModel event_log.Model,
	mongoStartSession mongo.StartSession,
	serviceName string,
) Consumer {
	return Consumer{
		logger:            logger,
		stanConn:          stan,
		fileModel:         fileModel,
		eventLogModel:     eventLogModel,
		mongoStartSession: mongoStartSession,
		serviceName:       serviceName,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
