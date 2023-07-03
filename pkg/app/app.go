package app

import (
	"github.com/msharbaji/grpc-go-example/internal/server"
	"github.com/rs/zerolog/log"
)

type App struct {
	grpcPort    string
	hmacSecrets map[string]string
	grpcServer  server.Grpc
}

func NewApp(grpcPort string, hmacSecrets map[string]string) (*App, error) {
	grpcServer, err := server.NewGrpcServer(grpcPort, hmacSecrets)
	if err != nil {
		return nil, err
	}

	return &App{
		grpcPort:    grpcPort,
		hmacSecrets: hmacSecrets,
		grpcServer:  *grpcServer,
	}, nil
}

func (a App) Run() error {
	go a.grpcServer.Start()

	a.grpcServer.HandleShutdown()
	// Initiate graceful shutdown
	log.Info().Msg("Received termination signal. Shutting down gRPC server...")
	if err := a.grpcServer.Stop(); err != nil {
		return err
	}

	return nil
}
