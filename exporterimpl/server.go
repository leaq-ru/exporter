package exporterimpl

import (
	"github.com/leaq-ru/exporter/cached_export"
	"github.com/leaq-ru/exporter/consumer"
	"github.com/leaq-ru/exporter/exporter_bucket"
	"github.com/leaq-ru/exporter/file"
	"github.com/leaq-ru/proto/codegen/go/exporter"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

type server struct {
	exporter.UnimplementedExporterServer
	logger            zerolog.Logger
	exporterBucket    exporter_bucket.ExporterBucket
	companyClient     parser.CompanyClient
	cityClient        parser.CityClient
	categoryClient    parser.CategoryClient
	fileModel         file.Model
	cachedExportModel cached_export.Model
	processAsync      consumer.ProcessAsync
}
