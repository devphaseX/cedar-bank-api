package gapi

import (
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/pb"
	"github.com/devphasex/cedar-bank-api/token"
	"github.com/devphasex/cedar-bank-api/util"
)

type GrpcServer struct {
	pb.UnimplementedSimpleBankServer
	tokenMaker token.Maker
	store      db.Store
	config     *util.Config
}

func NewGrpcServer(store db.Store, config *util.Config) (*GrpcServer, error) {

	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &GrpcServer{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
