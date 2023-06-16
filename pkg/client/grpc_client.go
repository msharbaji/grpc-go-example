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
	GetVersion(ctx context.Context, version string) (*pb.VersionResponse, error)

	// GetUser GetUser user by id, email or username.
	GetUser(ctx context.Context, identifier string, identifierType string) (*pb.User, error)

	// ListUsers ListUsers list all users.
	ListUsers(tx context.Context) (*pb.Users, error)

	// CreateUser create a new user.
	CreateUser(ctx context.Context, username string, email string) (*pb.User, error)

	// UpdateUser update a user.
	UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error)

	// DeleteUser delete a user.
	DeleteUser(ctx context.Context, identifier string, identifierType string) (*pb.User, error)
}

type client struct {
	pb.VersionServiceClient
	pb.UserServiceClient
}

func (c *client) ListUsers(ctx context.Context) (*pb.Users, error) {
	empty := &emptypb.Empty{}

	res, err := c.UserServiceClient.ListUsers(ctx, empty)
	if err != nil {
		log.Error().Err(err).Msg("failed to list users")
		return nil, err
	}
	return res, nil
}

func (c *client) GetVersion(ctx context.Context, version string) (*pb.VersionResponse, error) {

	versionReq := &pb.VersionRequest{
		Version: version,
	}

	res, err := c.VersionServiceClient.GetVersion(ctx, versionReq)
	if err != nil {
		log.Error().Err(err).Msg("failed to get version")
		return nil, err
	}
	return res, nil
}

func (c *client) GetUser(ctx context.Context, identifier string, identifierType string) (*pb.User, error) {
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

	res, err := c.UserServiceClient.GetUser(ctx, userReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) CreateUser(ctx context.Context, username string, email string) (*pb.User, error) {
	userReq := &pb.CreateUserRequest{
		Username: username,
		Email:    email,
	}

	res, err := c.UserServiceClient.CreateUser(ctx, userReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {

	res, err := c.UserServiceClient.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) DeleteUser(ctx context.Context, identifier string, identifierType string) (*pb.User, error) {

	userReq := &pb.DeleteUserRequest{}

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
	res, err := c.UserServiceClient.DeleteUser(ctx, userReq)

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
