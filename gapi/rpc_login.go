package gapi

import (
	"context"
	"errors"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/pb"
	"github.com/devphasex/cedar-bank-api/util/hash"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrMismatchCredential = errors.New("Invalid credential email or password mismatch")

func (s *GrpcServer) SigninUser(ctx context.Context, req *pb.CreateSigninRequest) (*pb.CreateSigninResponse, error) {

	ag := hash.DefaultArgonHash()

	user, err := s.store.GetUserByUniqueID(ctx, db.GetUserByUniqueIDParams{
		Email: pgtype.Text{
			String: req.ID,
			Valid:  true,
		},

		Username: pgtype.Text{
			String: req.ID,
			Valid:  true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, ErrMismatchCredential.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to sign in: %s", err.Error())
	}

	passwordHashByte, passwordSaltByte := hash.ArgonStringDecode(user.HashedPassword, user.PasswordSalt)

	if err = ag.Compare(passwordHashByte, passwordSaltByte, []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrMismatchCredential.Error())
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.AccessTokenTime)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.RefreshTokenTime)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:        pgtype.UUID{Bytes: [16]byte(refreshPayload.ID), Valid: true},
		OwnerID:   user.ID,
		UserAgent: "",
		ClientIp: pgtype.Text{
			String: "",
			Valid:  true,
		},
		RefreshToken: refreshToken,
		ExpiredAt: pgtype.Timestamptz{
			Time:  refreshPayload.ExpiresAt.Time,
			Valid: true,
		},
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rsp := &pb.CreateSigninResponse{
		SessionID:             uuid.UUID(session.ID.Bytes).String(),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
		User:                  convertDbUser(user),
	}

	return rsp, nil
}
