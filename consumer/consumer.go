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

type state struct {
	sub               stan.Subscription
	subscribeCalledOK bool
	drain             bool
	done              chan struct{}
}

type Consumer struct {
	logger            zerolog.Logger
	stanConn          stan.Conn
	exporterBucket    exporter_bucket.ExporterBucket
	companyClient     parser.CompanyClient
	fileModel         file.Model
	rowModel          row.Model
	cachedExportModel cached_export.Model
	mongoStartSession mongo.StartSession
	serviceName       string
	state             *state
}
