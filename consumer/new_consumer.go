package consumer

import (
	"github.com/leaq-ru/exporter/cached_export"
	"github.com/leaq-ru/exporter/exporter_bucket"
	"github.com/leaq-ru/exporter/file"
	"github.com/leaq-ru/exporter/processing_export"
	"github.com/leaq-ru/exporter/row"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"github.com/nats-io/stan.go"
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
