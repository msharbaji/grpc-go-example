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

func (s *versionServiceServer) GetVersion(_ context.Context, req *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	response := &pb.GetVersionResponse{
		Version: req.GetVersion(),
	}
	return response, nil
}
