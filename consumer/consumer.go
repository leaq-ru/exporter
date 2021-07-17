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

type state struct {
	sub               stan.Subscription
	subscribeCalledOK bool
	drain             bool
	done              chan struct{}
}

type Consumer struct {
	logger                zerolog.Logger
	stanConn              stan.Conn
	exporterBucket        exporter_bucket.ExporterBucket
	companyClient         parser.CompanyClient
	fileModel             file.Model
	rowModel              row.Model
	cachedExportModel     cached_export.Model
	processingExportModel processing_export.Model
	serviceName           string
	state                 *state
}
