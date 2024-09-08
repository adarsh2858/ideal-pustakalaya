package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	grpc "github.com/loupe-co/{{.repoName}}/cmd/server/grpc"
	pubsub "github.com/loupe-co/{{.repoName}}/cmd/server/pubsub"
	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/{{.repoName}}/internal/handlers"
	"github.com/loupe-co/go-common"
	configUtil "github.com/loupe-co/go-common/config"
	"github.com/loupe-co/go-common/errors"
	"github.com/loupe-co/go-loupe-logger/log"
)

func main() {
	// Get the service config
	cfg := config.Config{}
	err := configUtil.Load(
		&cfg,
		configUtil.FromENV(),
		configUtil.FromLocalYAML("/var/app-secrets/credentials_sentry.yaml"),
		configUtil.SetDefaultENV("project", "local"),
		configUtil.SetExportENVFromConfig(true),
	)
	if err != nil {
		// use fmt cause log hasn't been initialized yet (logger requires certain env vars ensured first)
		fmt.Println(errors.Wrap(err, "error loading {{.repoName}} config from env"))
		return
	}

	// Initialize logger
	l := log.InitLogger()
	defer l.Close()

	// Create cancel-able context
	ctx, cancel := context.WithCancel(context.Background())

	// Get the service handlers
	handles := handlers.New(cfg)

	// Get the grpc server
	grpcServer := grpc.New(cfg, handles)

	// Get the pubsub server
	pubsubServer, err := pubsub.New(ctx, cfg, handles)
	if err != nil {
		log.Error(errors.Wrap(err, "error getting pubsub server"))
		return
	}

	// Setup the os signal server
	sigServer := common.NewSigServer()
	sigServer.Handle(func(sig os.Signal) error {
		log.Infof("handling os signal %s", sig.String())
		return nil
	}, os.Interrupt, syscall.SIGTERM)

	// Start all servers and wait for any to error
	log.Info("Server starting")
	err = common.ServerListenMux(sigServer, grpcServer, pubsubServer)

	// Gracefully cancel/cleanup resources
	cancel()

	log.Info("Server exiting")
}
