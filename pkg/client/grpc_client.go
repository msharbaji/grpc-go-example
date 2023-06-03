package client

import (
	"context"
	"fmt"
	"github.com/msharbaji/grpc-go-example/api/pb"
	"github.com/msharbaji/grpc-go-example/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type Client interface {
	// GetVersion GetVersion application version.
	GetVersion(version string) (*pb.VersionResponse, error)

	// GetUser GetUser user by id, email or username.
	GetUser(identifier string, identifierType string) (*pb.User, error)

	// ListUsers ListUsers list all users.
	ListUsers() (*pb.Users, error)
}

type client struct {
	pb.VersionServiceClient
	pb.UserServiceClient
}

func (c *client) ListUsers() (*pb.Users, error) {
	empty := &emptypb.Empty{}

	res, err := c.UserServiceClient.ListUsers(context.Background(), empty)
	if err != nil {
		log.Error().Err(err).Msg("failed to list users")
		return nil, err
	}
	return res, nil
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

func (c *client) GetUser(identifier string, identifierType string) (*pb.User, error) {
	userReq := &pb.GetUserRequest{}

	switch strings.ToLower(identifierType) {
	case "id":
		userReq.Id = &identifier
	case "email":
		userReq.Email = &identifier
	case "username":
		userReq.Username = &identifier
	default:
		return nil, fmt.Errorf("invalid identifier type: %s", identifierType)
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
