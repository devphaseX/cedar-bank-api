package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://postgres:password@localhost:5432/cedar-bank?sslmode=disable"
)

var (
	testQueries Store
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var err error

	config, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	// Set some reasonable pool limits
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Set some reasonable timeouts
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.ConnConfig.RuntimeParams["statement_timeout"] = "30000" // 30 seconds

	testDB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("connection to db failed: %v", err)
	}
	defer testDB.Close()

	testQueries = NewStore(testDB)

	os.Exit(m.Run())
}
