package exporter_async

import (
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/event_log"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	minioClient *minio.Client,
	companyClient parser.CompanyClient,
	fileModel file.Model,
	eventLogModel event_log.Model,
	mongoStartSession mongo.StartSession,
	serviceName string,
) Consumer {
	return Consumer{
		logger:            logger,
		stanConn:          stanConn,
		minioClient:       minioClient,
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
