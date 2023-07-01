package handlers

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/msharbaji/grpc-go-example/pkg/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (s *userServiceServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, ok := users[req.GetUsername()]
	if ok {
		log.Error().Msg("user already exists")
		return nil, nil
	}

	user := &pb.User{
		Id:       uuid.New().String(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	users[req.GetUsername()] = user
	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser gets a user
func (s *userServiceServer) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Check if the request contains an ID, email, or username
	if id := req.GetId(); id != "" {
		for _, user := range users {
			if user.GetId() == id {
				return &pb.GetUserResponse{User: user}, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %s", id)
	}

	if email := req.GetEmail(); email != "" {
		// Retrieve user by email
		for _, user := range users {
			if user.GetEmail() == email {
				return &pb.GetUserResponse{User: user}, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "user not found with email: %s", email)
	}

	if username := req.GetUsername(); username != "" {
		// Retrieve user by username
		for _, user := range users {
			if user.GetUsername() == username {
				return &pb.GetUserResponse{User: user}, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "user not found with username: %s", username)
	}

	return nil, status.Error(codes.InvalidArgument, "missing ID, email, or username in the request")
}

// UpdateUser updates a user
func (s *userServiceServer) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	username := req.GetUsername()
	user, ok := users[username]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", username)
	}

	if email := req.GetEmail(); email != "" {
		user.Email = email
	}

	user.UpdatedAt = timestamppb.Now()

	users[username] = user

	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

// DeleteUser deletes a user
func (s *userServiceServer) DeleteUser(_ context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user := users[req.GetId()]
	delete(users, req.GetId())
	log.Info().Msgf("user deleted %s", req.GetId())
	return &pb.DeleteUserResponse{
		User: user,
	}, nil

}

// ListUsers lists all users
func (s *userServiceServer) ListUsers(_ context.Context, _ *emptypb.Empty) (*pb.ListUsersResponse, error) {
	var usersList []*pb.User
	for _, user := range users {
		usersList = append(usersList, user)
	}

	return &pb.ListUsersResponse{
		Users: usersList,
	}, nil
}
