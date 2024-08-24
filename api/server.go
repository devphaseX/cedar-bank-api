package api

import (
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	config *util.Config
}

func NewServer(store db.Store, config *util.Config) *Server {
	server := &Server{
		store:  store,
		config: config,
	}

	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountByID)
	router.GET("/accounts", server.getAccountList)
	router.POST("/transfer", server.transferTx)
	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
