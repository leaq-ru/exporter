package exporterimpl

import (
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/consumer"
	"github.com/nnqq/scr-exporter/exporter_bucket"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

type server struct {
	exporter.UnimplementedExporterServer
	logger            zerolog.Logger
	exporterBucket    exporter_bucket.ExporterBucket
	companyClient     parser.CompanyClient
	fileModel         file.Model
	cachedExportModel cached_export.Model
	processAsync      consumer.ProcessAsync
	mongoStartSession mongo.StartSession
}
