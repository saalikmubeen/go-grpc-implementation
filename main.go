package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	// Our Migration Source is the file system (file://db/migrations)
	// where our migration files are stored.
	// file:///absolute/path   |   file://relative/path
	// Other sources include: gitlab, s3, github, etc.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"

	"github.com/saalikmubeen/backend-masterclass-go/api"
	generated_db "github.com/saalikmubeen/backend-masterclass-go/db/sqlc"
	"github.com/saalikmubeen/backend-masterclass-go/grpcApi"
	"github.com/saalikmubeen/backend-masterclass-go/pb"
	"github.com/saalikmubeen/backend-masterclass-go/utils"
)

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)

	}

	// fmt.Println(config)

	conn, err := sql.Open(config.DBDriver, config.DBURI)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Another way of running database migrations
	// (directly from the GO code)
	runDBMigration(config.DBMigrationsURL, config.DBURI)

	store := generated_db.NewStore(conn)

	// Start the HTTP server
	// startHTTPGinServer(config, store)

	// Run the http Gateway server on separate Go routine.
	go startHTTPGatewayServer(config, store)

	// Start the gRPC server on main Go routine.
	startGRPCServer(config, store)
}

// ** HTTP Server (Serves the HTTP requests)
func startHTTPGinServer(config utils.Config, store generated_db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	fmt.Printf("Starting HTTP Gin server at %s....!\n", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}

// GRPC Server  (Serves the GRPC calls)
func startGRPCServer(config utils.Config, store generated_db.Store) {
	server, err := grpcApi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// ** create a new gRPC server
	grpcServer := grpc.NewServer()

	// register the service with the gRPC server
	pb.RegisterSimpleBankServiceServer(grpcServer, server)

	// *** Reflection
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen for GRPC Server: %v", err)
	}

	fmt.Printf("Starting GRPC server at %s....!\n", listener.Addr().String())
	// start the server
	// bind the port and listener to the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start GRPC server: %v", err)
	}

}

// ** HTTP Gateway Server (serves the HTTP requests through the GRPC)
func startHTTPGatewayServer(config utils.Config, store generated_db.Store) {

	server, err := grpcApi.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// The protocol buffer compiler generates camelCase JSON tags
	// that are used by default.
	// Inorder to use the exact case used in the proto files:
	// To use the same case (like camel case) names  for the JSON fields
	// as the proto field names, we need to use the UseProtoNames option
	// in the JSONPb struct.
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("Failed to register simple bank  service handler server: %v", err)
	}

	// This mux will receive all HTTP requests from the client
	httpMux := http.NewServeMux()

	// Register the gRPC Gateway mux with the HTTP mux
	// To convert the HTTP requests to gRPC requests, route them
	// to the grpcMux.
	httpMux.Handle("/", grpcMux)

	// http gateway listener
	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Fatalf("Failed to listen for HHTP Gateway Server: %v", err)
	}

	log.Printf("STARTING HTTP GATEWAY SERVER ON PORT %s....!", listener.Addr().String())
	// start the server
	// bind the port and listener to the gRPC server
	if err := http.Serve(listener, httpMux); err != nil {
		log.Fatalf("Cannot start HTTP Gateway Server...: %v", err)
	}
}
