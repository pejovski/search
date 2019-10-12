package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/search/controller"
	"github.com/pejovski/search/factory"
	"github.com/pejovski/search/repository"
	httpServer "github.com/pejovski/search/server/http"
)

const (
	serverShutdownTimeout = 3 * time.Second
)

func main() {
	esClient := factory.CreateESClient(fmt.Sprintf(
		"http://%s:%s",
		os.Getenv("ES_HOST"),
		os.Getenv("ES_PORT"),
	))

	searchRepository := repository.NewESProductRepository(esClient)
	searchController := controller.NewSearch(searchRepository)

	serverHandler := httpServer.NewHandler(searchController)
	serverRouter := httpServer.NewRouter(serverHandler)

	server := factory.CreateHttpServer(serverRouter, fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf(err.Error())
		}
	}()
	logrus.Infof("Server started at port: %s", os.Getenv("APP_PORT"))

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Receive shutdown signals.
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorln("Server shutdown failed", err)
	}
	logrus.Println("Server exited properly")
}
