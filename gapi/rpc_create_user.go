package gapi

import (
	"context"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/pb"
	"github.com/devphasex/cedar-bank-api/util/hash"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ag := hash.DefaultArgonHash()

	passwordHash, err := ag.GenerateHash([]byte(req.GetPassword()), nil)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	passwordHashStr, passwordSaltStr := hash.ArgonStringEncode(passwordHash)

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		Fullname:       req.GetFullname(),
		HashedPassword: passwordHashStr,
		PasswordSalt:   passwordSaltStr,
	}
	user, err := s.store.CreateUser(ctx, arg)

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.ConstraintName {
			case "users_username_key":
				return nil, status.Error(codes.AlreadyExists, "username already taken")

			case "users_email_key":
				return nil, status.Error(codes.AlreadyExists, "email already taken")
			}

		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertDbUser(user),
	}

	return rsp, nil
}
