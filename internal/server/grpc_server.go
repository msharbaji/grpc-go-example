package server

import (
	"errors"
	"fmt"
	handlers2 "github.com/msharbaji/grpc-go-example/internal/handlers"
	"github.com/msharbaji/grpc-go-example/pkg/middleware"
	"github.com/msharbaji/grpc-go-example/pkg/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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
func NewGrpcServer(port string, secrets map[string]string) (*Grpc, error) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.NewServerAuthInterceptor(secrets)),
		grpc.Creds(insecure.NewCredentials()),
	}
	s := &Grpc{
		address: fmt.Sprintf(":%s", port),
		server:  grpc.NewServer(opts...),
	}

	pb.RegisterVersionServiceServer(s.server, handlers2.NewVersionServiceServer())
	pb.RegisterUserServiceServer(s.server, handlers2.NewUserServiceServer())

	reflection.Register(s.server)
	return s, nil
}

// Start starts the grpc server
func (s *Grpc) Start() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	if err := s.server.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal().Err(err).Msg("failed to start gRPC server")
	}

	log.Info().Msgf("gRPC server started on %s", s.address)

}

// Stop stops the grpc server
func (s *Grpc) Stop() error {
	s.server.GracefulStop()
	return nil
}

func (s *Grpc) HandleShutdown() {
	// Create a channel to receive termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a termination signal is received
	<-signalCh

	// Initiate graceful shutdown
	log.Info().Msg("Received termination signal. Shutting down gRPC server...")
	if err := s.Stop(); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown gRPC server")
	}
}
