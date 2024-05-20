package grpcApi

import (
	generated_db "github.com/saalikmubeen/go-grpc-implementation/db/sqlc"
	"github.com/saalikmubeen/go-grpc-implementation/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func transformUser(user generated_db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
