package app

import (
	"errors"
	"github.com/msharbaji/grpc-go-example/api/pb"
	"github.com/msharbaji/grpc-go-example/internal/app/handlers"
	"github.com/msharbaji/grpc-go-example/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Grpc struct {
	address string
	server  *grpc.Server
}

// NewGrpcServer creates a new grpc server
func NewGrpcServer(address string, secrets map[string]string) (*Grpc, error) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.NewServerAuthInterceptor(secrets)),
		grpc.Creds(insecure.NewCredentials()),
	}
	s := &Grpc{
		address: address,
		server:  grpc.NewServer(opts...),
	}

	pb.RegisterVersionServiceServer(s.server, handlers.NewVersionServiceServer())
	pb.RegisterUserServiceServer(s.server, handlers.NewUserServiceServer())
	return s, nil
}

// Start starts the grpc server
func (s *Grpc) Start() error {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	go func() {
		if err := s.server.Serve(listen); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal().Err(err).Msg("failed to start gRPC server")
		}
	}()

	log.Info().Msgf("gRPC server started on %s", s.address)

	return nil
}

// Stop stops the grpc server
func (s *Grpc) Stop() error {
	s.server.GracefulStop()
	return nil
}

// Run starts the gRPC server and handles graceful shutdown
func (s *Grpc) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Create a channel to receive termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a termination signal is received
	<-signalCh

	// Initiate graceful shutdown
	log.Info().Msg("Received termination signal. Shutting down gRPC server...")
	if err := s.Stop(); err != nil {
		return err
	}

	return nil
}
