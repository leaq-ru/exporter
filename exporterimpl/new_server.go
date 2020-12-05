package exporterimpl

import (
	"github.com/nnqq/scr-exporter/exporter_async"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/rs/zerolog"
)

func NewServer(
	logger zerolog.Logger,
	fileModel file.Model,
	processAsync exporter_async.ProcessAsync,
	mongoStartSession mongo.StartSession,
) *server {
	return &server{
		logger:            logger,
		fileModel:         fileModel,
		processAsync:      processAsync,
		mongoStartSession: mongoStartSession,
	}
}
