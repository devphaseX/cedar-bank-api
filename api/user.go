package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/util/hash"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"min=2,required"`
	Email    string `json:"email" binding:"email,required"`
	Fullname string `json:"fullname" binding:"min=2,required"`
	Password string `json:"password" binding:"min=8,required"`
}

type userResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Fullname          string    `json:"fullname"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Fullname:          user.Fullname,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.Time,
		CreatedAt:         user.CreatedAt.Time,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, prettyValidateError(err))
		return
	}

	ag := hash.DefaultArgonHash()

	passwordHash, err := ag.GenerateHash([]byte(req.Password), nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.New("an error occurred while creating your account"))
		return
	}

	passwordHashStr, passwordSaltStr := hash.ArgonStringEncode(passwordHash)

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		Fullname:       req.Fullname,
		HashedPassword: passwordHashStr,
		PasswordSalt:   passwordSaltStr,
	}
	user, err := s.store.CreateUser(ctx, arg)

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.ConstraintName {
			case "users_username_key":
				ctx.JSON(http.StatusConflict, errors.New(fmt.Sprintf("username already taken")))
				return
			case "users_email_key":
				ctx.JSON(http.StatusConflict, errors.New(fmt.Sprintf("email already taken")))
				return
			}

		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)

	ctx.JSON(http.StatusCreated, sucessResponse(resp, "user created successfully"))
	return
}

type SigninRequest struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signinResponse struct {
	SessionID             string       `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) signin(ctx *gin.Context) {

	var req SigninRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, prettyValidateError(err))
		return
	}

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
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid credential email or password mismatch")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	passwordHashByte, passwordSaltByte := hash.ArgonStringDecode(user.HashedPassword, user.PasswordSalt)

	if err = ag.Compare(passwordHashByte, passwordSaltByte, []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid credential email or password mismatch")))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.AccessTokenTime)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.RefreshTokenTime)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	response := signinResponse{
		SessionID:             uuid.UUID(session.ID.Bytes).String(),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiresAt.Time,
		User:                  resp,
	}

	ctx.JSON(http.StatusOK, sucessResponse(response, "account signin successfully"))
}
