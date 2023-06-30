package handlers

import (
	"context"
	"github.com/msharbaji/grpc-go-example/pkg/pb"
)

type versionServiceServer struct {
	pb.UnimplementedVersionServiceServer
}

func NewVersionServiceServer() pb.VersionServiceServer {
	return &versionServiceServer{}
}

func (s *versionServiceServer) GetVersion(_ context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	response := &pb.VersionResponse{
		Version: req.GetVersion(),
	}
	return response, nil
}
