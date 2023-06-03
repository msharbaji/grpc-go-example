package handlers

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/msharbaji/grpc-go-example/api/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var users = map[string]*pb.User{
	"someone": {
		Id:       "1",
		Username: "someone",
		Email:    "someone@someone.com",
	},
	"someone_else": {
		Id:       "2",
		Username: "someone_else",
		Email:    "someonce2@someone.com",
		CreatedAt: &timestamp.Timestamp{
			Seconds: 1612345678,
		},
		UpdatedAt: &timestamp.Timestamp{
			Seconds: 1612345678,
		},
	},
}

// UserServiceServer is the user service server
type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// NewUserServiceServer creates a new user service server
func NewUserServiceServer() pb.UserServiceServer {
	return &userServiceServer{}
}

// CreateUser creates a new user
func (s *userServiceServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	_, ok := users[req.GetUsername()]
	if ok {
		log.Error().Msg("user already exists")
		return nil, nil
	}

	user := &pb.User{
		Id:       "1",
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	users[req.GetUsername()] = user
	return user, nil
}

// GetUser gets a user
func (s *userServiceServer) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// Check if the request contains an ID, email, or username
	if id := req.GetId(); id != "" {
		// Retrieve user by ID
		user, ok := users[id]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "user not found with ID: %s", id)
		}
		return user, nil
	}

	if email := req.GetEmail(); email != "" {
		// Retrieve user by email
		for _, user := range users {
			if user.GetEmail() == email {
				return user, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "user not found with email: %s", email)
	}

	if username := req.GetUsername(); username != "" {
		// Retrieve user by username
		for _, user := range users {
			if user.GetUsername() == username {
				return user, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "user not found with username: %s", username)
	}

	return nil, status.Error(codes.InvalidArgument, "missing ID, email, or username in the request")
}

// UpdateUser updates a user
func (s *userServiceServer) UpdateUser(_ context.Context, req *pb.User) (*pb.User, error) {
	user, ok := users[req.GetUsername()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user not found: %s", req.GetUsername()))
	}

	users[req.GetUsername()] = req
	return user, nil
}

// DeleteUser deletes a user
func (s *userServiceServer) DeleteUser(_ context.Context, req *pb.DeleteUserRequest) (*pb.User, error) {
	delete(users, req.GetId())
	log.Info().Msgf("user deleted %s", req.GetId())
	return &pb.User{}, nil

}

// ListUsers lists all users
func (s *userServiceServer) ListUsers(context.Context, *emptypb.Empty) (*pb.Users, error) {
	var usersList []*pb.User
	for _, user := range users {
		usersList = append(usersList, user)
	}

	return &pb.Users{
		Users: usersList,
	}, nil
}
