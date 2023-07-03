package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/msharbaji/grpc-go-example/pkg/app"
	"github.com/rs/zerolog/log"
)

const (
	//nolint:unused
	version = "local"
)

var (
	grpcPort    = kingpin.Flag("grpc-port", "gRPC port").Envar("GRPC_PORT").Default("50051").String()
	hmacSecrets = kingpin.Flag("hmac-secrets", "Key-value pair for secret").Envar("HMAC_SECRETS").Default("my-secret-key=my-secret-value").StringMap()
)

func main() {
	// parse command line flags
	kingpin.Parse()

	log.Info().Str("AppVersion", version).Msg("starting api")

	_app, err := app.NewApp(*grpcPort, *hmacSecrets)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}

	if err := _app.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}

}
