package api

import (
	"errors"
	"net/http"

	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

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

	tx, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		if errors.As(err, &db.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		if errors.As(err, &db.ErrFundNotSufficient) {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, sucessResponse(tx))
}
