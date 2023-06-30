package main

import (
	"github.com/msharbaji/grpc-go-example/internal/utils"
	"github.com/msharbaji/grpc-go-example/pkg/server"
	"github.com/rs/zerolog/log"
)

const (
	//nolint:unused
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
	grpcServer, err := server.NewGrpcServer(cfg.Endpoint, secrets)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}

	//defer wg.Done()
	if err := grpcServer.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run grpc server")
	}

}
