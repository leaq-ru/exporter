package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/exporter_bucket"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/processing_export"
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
	processingExportModel processing_export.Model,
	serviceName string,
) Consumer {
	return Consumer{
		logger:                logger,
		stanConn:              stanConn,
		exporterBucket:        exporterBucket,
		companyClient:         companyClient,
		fileModel:             fileModel,
		rowModel:              rowModel,
		cachedExportModel:     cachedExportModel,
		processingExportModel: processingExportModel,
		serviceName:           serviceName,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
