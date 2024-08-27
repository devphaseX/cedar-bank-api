package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required"`
	ToAccountID   int64   `json:"to_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

var ErrFailed = errors.New("failed")

func (s *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, prettyValidateError(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	if !s.validateAccount(ctx, arg.FromAccountID, req.Currency, true) {
		return
	}

	if !s.validateAccount(ctx, arg.ToAccountID, req.Currency, false) {
		return
	}

	tx, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		if errors.Is(err, db.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		if errors.Is(err, db.ErrFundNotSufficient) {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sucessResponse(tx))
}

func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string, ownerCheck bool) bool {
	account, err := s.store.GetAccountByID(ctx, accountID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("account [%d] not found", accountID)))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if ownerCheck {
		authUser := Auth(ctx)

		if authUser.UserId != account.OwnerID {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authorized")))
			return false
		}
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusNotFound, errorResponse(
			fmt.Errorf("account [%d] currency mismatch: expected %v but got %v",
				accountID, account.Currency, currency)),
		)
		return false
	}

	return true
}
