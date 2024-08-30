package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/devphasex/cedar-bank-api/api"
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	_ "github.com/devphasex/cedar-bank-api/doc/statik"
	"github.com/devphasex/cedar-bank-api/gapi"
	"github.com/devphasex/cedar-bank-api/pb"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	store := db.NewStore(conn)
	go runGrpcServer(store, config)
	runGrpcGatewayServer(store, config)
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
		log.Fatal("cannot create listerner:", err)
	}

	log.Printf("start gRPC server at %s", ln.Addr().String())
	err = grpcServer.Serve(ln)

	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGrpcGatewayServer(store db.Store, config *util.Config) {
	var err error
	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Assuming your gRPC server is running on a different port, e.g., ":9090"
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // gRPC connection options

	err = pb.RegisterSimpleBankHandlerFromEndpoint(ctx, grpcMux, config.GrpcServerAddress, opts)
	if err != nil {
		log.Fatal("cannot register handler from endpoint:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFs, err := fs.New()

	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	// Add CORS middleware if necessary
	// handler := cors.Default().Handler(mux)

	ln, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	defer ln.Close()

	log.Printf("starting HTTP server at %s", ln.Addr().String())
	err = http.Serve(ln, mux)
	if err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}
