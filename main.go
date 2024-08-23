package main

import (
	"context"
	"log"
	"time"

	"github.com/devphasex/cedar-bank-api/api"
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/jackc/pgx/v5/pgxpool"
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

	server := api.NewServer(db.NewStore(conn), config)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
