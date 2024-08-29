package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/devphasex/cedar-bank-api/api"
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/gapi"
	"github.com/devphasex/cedar-bank-api/pb"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var err error

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal(err)
	}

	pgConfig, err := pgxpool.ParseConfig(config.DbSource)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	// Set some reasonable pool limits
	pgConfig.MaxConns = 20
	pgConfig.MinConns = 2
	pgConfig.MaxConnLifetime = time.Hour
	pgConfig.MaxConnIdleTime = 30 * time.Minute

	// Set some reasonable timeouts
	pgConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	pgConfig.ConnConfig.RuntimeParams["statement_timeout"] = "30000" // 30 seconds

	conn, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatalf("connection to db failed: %v", err)
	}
	defer conn.Close()

	runGrpcServer(db.NewStore(conn), config)
}

func runGinServer(store db.Store, config *util.Config) {
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(config.HttpServerAddress); err != nil {
		log.Fatal(err)
	}
}

func runGrpcServer(store db.Store, config *util.Config) {
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	server, err := gapi.NewGrpcServer(store, config)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	ln, err := net.Listen("tcp", config.GrpcServerAddress)

	defer ln.Close()

	if err != nil {
		log.Fatal("cannot create listerner")
	}

	log.Printf("start gRPC server at %s", ln.Addr().String())
	err = grpcServer.Serve(ln)

	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}
