package exporterimpl

import (
	"github.com/nnqq/scr-exporter/consumer"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-exporter/store"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

type server struct {
	exporter.UnimplementedExporterServer
	logger            zerolog.Logger
	store             store.Store
	companyClient     parser.CompanyClient
	fileModel         file.Model
	processAsync      consumer.ProcessAsync
	mongoStartSession mongo.StartSession
}
