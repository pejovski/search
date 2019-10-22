package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/search/controller"
	"github.com/pejovski/search/factory"
	"github.com/pejovski/search/pkg/signals"
	"github.com/pejovski/search/repository/es"
	"github.com/pejovski/search/server/api"
)

const (
	shutdownDuration = 3 * time.Second
)

func main() {
	esClient := factory.CreateESClient(fmt.Sprintf(
		"http://%s:%s",
		os.Getenv("ES_HOST"),
		os.Getenv("ES_PORT"),
	))

	esRepo := es.NewRepository(esClient)
	c := controller.New(esRepo)

	ctx := signals.Context()

	serverAPI := api.NewServer(c)
	serverAPI.Run(ctx)

	logrus.Infof("allowing %s for graceful shutdown to complete", shutdownDuration)
	<-time.After(shutdownDuration)
}
