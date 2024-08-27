package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"oneof=currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUser := Auth(ctx)

	arg := db.CreateAccountParams{
		OwnerID:  authUser.UserId,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(ctx, arg)

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.ConstraintName {
			case "fk_accounts_users":
				ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user with ID %d not found", arg.OwnerID)))
				return
			case "unique_owner_currency":
				ctx.JSON(http.StatusConflict, errorResponse(fmt.Errorf("account already exists for user ID %d with currency %s", arg.OwnerID, arg.Currency)))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sucessResponse(account))
}

type GetAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"min=1"`
}

func (s *Server) getAccountByID(ctx *gin.Context) {
	var req GetAccountByIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccountByID(ctx, req.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("account with id '%v' not found", req.ID)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authUser := Auth(ctx)
	if account.OwnerID != authUser.UserId {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authorized")))
		return
	}

	ctx.JSON(http.StatusOK, sucessResponse(account))
}

type GetAccountList struct {
	Page    int32 `form:"page" binding:"required,gt=0"`
	PerPage int32 `form:"per_page" binding:"required,min=5"`
}

func (s *Server) getAccountList(ctx *gin.Context) {
	var req GetAccountList

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUser := Auth(ctx)

	arg := db.GetAccountsParams{
		Offset: int64((req.Page - 1) * req.PerPage),
		Limit:  int64(req.PerPage),
		UserID: pgtype.Int8{
			Valid: true,
			Int64: authUser.UserId,
		},
	}

	accounts, err := s.store.GetAccounts(ctx, arg)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sucessResponse(accounts))
}
