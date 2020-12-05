package main

import (
	"context"
	"github.com/nnqq/scr-exporter/call"
	"github.com/nnqq/scr-exporter/config"
	"github.com/nnqq/scr-exporter/event_log"
	"github.com/nnqq/scr-exporter/exporter_async"
	"github.com/nnqq/scr-exporter/exporterimpl"
	"github.com/nnqq/scr-exporter/file"
	"github.com/nnqq/scr-exporter/logger"
	"github.com/nnqq/scr-exporter/minio"
	"github.com/nnqq/scr-exporter/mongo"
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

	companyClient, err := call.NewClients(cfg.Service.Parser)
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
		cfg.S3.Region,
		cfg.S3.ExporterBucketName,
		cfg.S3.Secure,
	)
	logg.Must(err)

	consumer := exporter_async.NewConsumer(
		logg.ZL,
		stanConn,
		minioClient,
		companyClient,
		file.NewModel(db),
		event_log.NewModel(db),
		db.Client().StartSession,
		cfg.ServiceName,
	)
	logg.Must(consumer.Subscribe())

	grpcSrv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, health.NewServer())
	exporter.RegisterExporterServer(grpcSrv, exporterimpl.NewServer(
		logg.ZL,
		file.NewModel(db),
		consumer.ProcessAsync,
		db.Client().StartSession,
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
		graceful.HandleSignals(grpcSrv.GracefulStop, consumer.GracefulStop)
	}()
	go func() {
		defer wg.Done()
		logg.Must(grpcSrv.Serve(lis))
	}()
	go func() {
		defer wg.Done()
		logg.Must(consumer.Serve())
	}()
	wg.Wait()
}
