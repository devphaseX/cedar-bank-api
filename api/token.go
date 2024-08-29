package api

import (
	"errors"
	"net/http"
	"time"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, prettyValidateError(err))
		return
	}

	session, err := s.store.GetSessionByUniqueID(ctx, db.GetSessionByUniqueIDParams{
		RefreshToken: pgtype.Text{
			String: req.RefreshToken,
			Valid:  true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked.Bool {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("invalid session")))
		return
	}

	payload, err := s.tokenMaker.VerifyToken(session.RefreshToken)

	if payload.UserId != session.OwnerID {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("invalid session")))
		return
	}

	if err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session expired")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to refresh token")))
		return
	}

	if time.Now().After(session.ExpiredAt.Time) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session expired")))
		return
	}

	user, err := s.store.GetUserByUniqueID(ctx, db.GetUserByUniqueIDParams{
		ID: pgtype.Int8{
			Int64: payload.UserId,
			Valid: true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session expired")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	accessTokenStr, accessPayload, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.AccessTokenTime)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := renewAccessTokenResponse{
		AccessToken:          accessTokenStr,
		AccessTokenExpiredAt: accessPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, sucessResponse(response, "access token refreshed succesfully"))
}
