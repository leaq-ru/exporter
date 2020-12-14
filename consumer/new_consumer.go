package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/exporter_bucket"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-exporter/row"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	exporterBucket exporter_bucket.ExporterBucket,
	companyClient parser.CompanyClient,
	fileModel file.Model,
	rowModel row.Model,
	cachedExportModel cached_export.Model,
	mongoStartSession mongo.StartSession,
	serviceName string,
) Consumer {
	return Consumer{
		logger:            logger,
		stanConn:          stanConn,
		exporterBucket:    exporterBucket,
		companyClient:     companyClient,
		fileModel:         fileModel,
		rowModel:          rowModel,
		cachedExportModel: cachedExportModel,
		mongoStartSession: mongoStartSession,
		serviceName:       serviceName,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
