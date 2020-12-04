package exporterimpl

import (
	"github.com/nnqq/scr-exporter/exporter_async"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

type server struct {
	exporter.UnimplementedExporterServer
	logger            zerolog.Logger
	companyClient     parser.CompanyClient
	fileModel         file.Model
	processAsync      exporter_async.ProcessAsync
	mongoStartSession mongo.StartSession
}
