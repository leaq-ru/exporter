package exporterimpl

import (
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/consumer"
	"github.com/nnqq/scr-exporter/exporter_bucket"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

func NewServer(
	logger zerolog.Logger,
	exporterBucket exporter_bucket.ExporterBucket,
	companyClient parser.CompanyClient,
	fileModel file.Model,
	cachedExportModel cached_export.Model,
	processAsync consumer.ProcessAsync,
	mongoStartSession mongo.StartSession,
) *server {
	return &server{
		logger:            logger,
		exporterBucket:    exporterBucket,
		companyClient:     companyClient,
		fileModel:         fileModel,
		cachedExportModel: cachedExportModel,
		processAsync:      processAsync,
		mongoStartSession: mongoStartSession,
	}
}
