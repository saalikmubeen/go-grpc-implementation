package grpcApi

import (
	"context"
	"fmt"

	"github.com/saalikmubeen/go-grpc-implementation/pb"
)

func (server *server) HelloFromServer(ctx context.Context, req *pb.HelloFromServerRequest) (*pb.HelloFromServerResponse, error) {

	metadata := server.extractMetadata(ctx)

	return &pb.HelloFromServerResponse{
		Message:    "Hello from the server",
		ClientIp:   metadata.ClientIP,
		UserAgent:  metadata.UserAgent,
		ServerName: fmt.Sprintf("HTTP: %s, gRPC: %s", server.config.HTTPServerAddress, server.config.GRPCServerAddress),
	}, nil
}
