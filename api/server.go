package api

import (
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/token"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
	config     *util.Config
}

func NewServer(store db.Store, config *util.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidators(validator)
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/auth/sign-up", s.createUser)
	router.POST("/auth/sign-in", s.signin)

	authProtectedRoute := router.Group("/").Use(AuthMiddleware(s.tokenMaker))

	authProtectedRoute.POST("/accounts", s.createAccount)
	authProtectedRoute.GET("/accounts/:id", s.getAccountByID)
	authProtectedRoute.GET("/accounts", s.getAccountList)
	authProtectedRoute.POST("/transfer", s.createTransfer)

	s.router = router
}
