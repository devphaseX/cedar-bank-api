package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id,required"`
	ToAccountID   int64   `json:"to_account_id,required"`
	Amount        float64 `json:"amount,required"`
	Currency      string  `json:"currency,required" binding:"oneof=USD EUR CAD"`
}

var ErrFailed = errors.New("failed")

func (s *Server) transferTx(ctx *gin.Context) {
	var req TransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	if !s.validateAccount(ctx, arg.FromAccountID, req.Currency) {
		return
	}

	if !s.validateAccount(ctx, arg.ToAccountID, req.Currency) {
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

func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccountByID(ctx, accountID)

	if err != nil {
		if errors.Is(err, db.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
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
