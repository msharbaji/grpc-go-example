package client

import (
	"context"
	"github.com/msharbaji/grpc-go-example/api/pb"
	"github.com/msharbaji/grpc-go-example/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	// GetVersion GetVersion application version.
	GetVersion(version string) (*pb.VersionResponse, error)

	// GetUser GetVersion user by id.
	GetUser(id string) (*pb.User, error)
}

type client struct {
	pb.VersionServiceClient
	pb.UserServiceClient
}

func (c *client) GetVersion(version string) (*pb.VersionResponse, error) {

	versionReq := &pb.VersionRequest{
		Version: version,
	}

	res, err := c.VersionServiceClient.GetVersion(context.Background(), versionReq)
	if err != nil {
		log.Error().Err(err).Msg("failed to get version")
		return nil, err
	}
	return res, nil
}

func (c *client) GetUser(id string) (*pb.User, error) {

	userReq := &pb.GetUserRequest{
		Id: id,
	}

	res, err := c.UserServiceClient.GetUser(context.Background(), userReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NewClient creates a new grpc client
func NewClient(endpoint, hmacKeyID, hmacSecret string) (Client, error) {
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(middleware.NewClientAuthInterceptor(hmacKeyID, hmacSecret)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to dial server")
	}

	v := pb.NewVersionServiceClient(conn)
	c := pb.NewUserServiceClient(conn)
	return &client{
		VersionServiceClient: v,
		UserServiceClient:    c,
	}, nil
}
