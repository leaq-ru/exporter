package main

import (
	"context"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-exporter/call"
	"github.com/nnqq/scr-exporter/config"
	"github.com/nnqq/scr-exporter/consumer"
	"github.com/nnqq/scr-exporter/exporter_bucket"
	"github.com/nnqq/scr-exporter/exporterimpl"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/logger"
	"github.com/nnqq/scr-exporter/minio"
	"github.com/nnqq/scr-exporter/mongo"
	"github.com/nnqq/scr-exporter/processing_export"
	"github.com/nnqq/scr-exporter/row"
	"github.com/nnqq/scr-exporter/stan"
	graceful "github.com/nnqq/scr-lib-graceful"
	"github.com/nnqq/scr-proto/codegen/go/exporter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	companyClient, cityClient, categoryClient, err := call.NewClients(cfg.Service.Parser)
	logg.Must(err)

	stanConn, err := stan.NewConn(cfg.ServiceName, cfg.STAN.ClusterID, cfg.NATS.URL)
	logg.Must(err)

	db, err := mongo.NewConn(ctx, cfg.ServiceName, cfg.MongoDB.URL)
	logg.Must(err)

	minioClient, err := minio.NewClient(
		ctx,
		cfg.S3.AccessKeyID,
		cfg.S3.SecretAccessKey,
		cfg.S3.Endpoint,
		cfg.S3.Secure,
	)
	logg.Must(err)

	exporterBucket, err := exporter_bucket.NewExporterBucket(
		ctx,
		minioClient,
		cfg.S3.ExporterBucketName,
		cfg.S3.Region,
	)
	logg.Must(err)

	fileModel := file.NewModel(db)
	cachedExportModel := cached_export.NewModel(db)
	processingExportModel := processing_export.NewModel(db)
	rowModel := row.NewModel(db)

	cons := consumer.NewConsumer(
		logg.ZL,
		stanConn,
		exporterBucket,
		companyClient,
		fileModel,
		rowModel,
		cachedExportModel,
		processingExportModel,
		cfg.ServiceName,
	)
	logg.Must(cons.Subscribe())

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	exporter.RegisterExporterServer(srv, exporterimpl.NewServer(
		logg.ZL,
		exporterBucket,
		companyClient,
		cityClient,
		categoryClient,
		fileModel,
		cachedExportModel,
		cons.ProcessAsync,
	))

	lis, err := net.Listen("tcp", strings.Join([]string{
		"0.0.0.0",
		cfg.Grpc.Port,
	}, ":"))
	logg.Must(err)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		graceful.HandleSignals(srv.GracefulStop, cons.GracefulStop)
	}()
	go func() {
		defer wg.Done()
		logg.Must(srv.Serve(lis))
	}()
	go func() {
		defer wg.Done()
		logg.Must(cons.Serve())
	}()
	wg.Wait()
}
