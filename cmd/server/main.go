package main

import (
	"github.com/msharbaji/grpc-go-example/internal/app"
	"github.com/msharbaji/grpc-go-example/internal/utils"
	"github.com/rs/zerolog/log"
)

const (
	version = "local"
)

func main() {

	// Load config
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Create a new grpc server
	secrets := map[string]string{
		cfg.KeyID: cfg.SecretKey,
	}
	grpcServer, err := app.NewGrpcServer(cfg.Endpoint, secrets)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}

	//defer wg.Done()
	if err := grpcServer.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run grpc server")
	}

}
