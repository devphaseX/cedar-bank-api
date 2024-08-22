package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://postgres:password@localhost:5432/cedar-bank?sslmode=disable"
)

var testQueries *Store

func TestMain(m *testing.M) {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("connection to db failed: %v", err)
	}
	defer conn.Close()

	testQueries = NewStore(conn)

	os.Exit(m.Run())
}
