package api

import (
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidators(validator)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountByID)
	router.GET("/accounts", server.getAccountList)
	router.POST("/transfer", server.createTransfer)
	router.POST("/auth/sign-up", server.createUser)
	router.POST("/auth/sign-in", server.signin)
	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
