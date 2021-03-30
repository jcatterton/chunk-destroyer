package main

import (
	"chunk-destroyer/config"
	"chunk-destroyer/pkg/producer"
	"chunk-destroyer/pkg/service"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.InitializeConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize config")
		os.Exit(1)
	}

	p, err := producer.CreateProducer(conf.KafkaBroker, conf.KafkaTopic)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create producer")
		os.Exit(1)
	}

	ctx := context.Background()
	s, err := service.InitializeService(ctx, conf)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize service")
		p.Produce("initialization_failed", "Failed to initialize service", true)
		os.Exit(1)
	}

	deletedChunks, err := s.Run(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to run service")
		p.Produce("initialization_failed", "Error running service", true)
		os.Exit(1)
	}

	logrus.WithField("deleted_chunks", deletedChunks).Info("Service ran successfully, terminating")
	p.Produce("service_complete", fmt.Sprintf("Service ran successfully, %v chunks deleted", deletedChunks), false)
	time.Sleep(30 * time.Second)
	os.Exit(0)
}
