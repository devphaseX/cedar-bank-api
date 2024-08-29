package gapi

import (
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertDbUser(user db.User) *pb.User {
	return &pb.User{
		Id:                int32(user.ID),
		Username:          user.Username,
		Email:             user.Email,
		Fullname:          user.Fullname,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}
}
