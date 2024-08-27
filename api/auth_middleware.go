package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/devphasex/cedar-bank-api/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayload    = "auth_payload"
	authorizationTypeBearer = "bearer"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		fmt.Println("header", authorizationHeader)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization token not set in header")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != "bearer" {
			err := errors.New(fmt.Sprintf("unsupported authorization type %s", authorizationType))

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayload, payload)

		ctx.Next()
	}
}

func Auth(ctx *gin.Context) *token.Payload {
	payload, ok := ctx.MustGet(authorizationPayload).(*token.Payload)

	if !ok {
		log.Fatal("payload not conform to the *Payload type")
		return nil
	}

	return payload
}
