package exporterimpl

import (
	"github.com/nnqq/scr-exporter/exporter_async"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/rs/zerolog"
)

func NewServer(
	logger zerolog.Logger,
	fileModel file.Model,
	companyClient parser.CompanyClient,
	processAsync exporter_async.ProcessAsync,
	mongoStartSession mongo.StartSession,
) *server {
	return &server{
		logger:            logger,
		fileModel:         fileModel,
		companyClient:     companyClient,
		processAsync:      processAsync,
		mongoStartSession: mongoStartSession,
	}
}
