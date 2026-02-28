// Command server runs the gRPC and HTTP API for the scheduler.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
	"scheduler/backend/internal/broadcast"
	"scheduler/backend/internal/config"
	"scheduler/backend/internal/gateway"
	"scheduler/backend/internal/server"
	"scheduler/backend/internal/service"
	"scheduler/backend/internal/store"
)

func main() {
	configuration := config.Load()
	flag.StringVar(&configuration.GRPCAddr, "grpc", configuration.GRPCAddr, "gRPC listen address")
	flag.StringVar(&configuration.HTTPAddr, "http", configuration.HTTPAddr, "HTTP gateway listen address")
	flag.StringVar(&configuration.DBPath, "db", configuration.DBPath, "SQLite database path")
	flag.Parse()

	appointmentStore, err := store.New(configuration.DBPath)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer appointmentStore.Close()

	broadcaster := broadcast.New()
	schedulerService := service.New(appointmentStore, broadcaster)
	grpcServer := server.New(schedulerService, broadcaster)
	httpGateway := gateway.New(schedulerService, broadcaster)

	listener, err := net.Listen("tcp", configuration.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	grpcServerInstance := grpc.NewServer()
	schedulerv1.RegisterSchedulerServiceServer(grpcServerInstance, grpcServer)
	reflection.Register(grpcServerInstance)
	go func() {
		log.Printf("gRPC listening on %v", configuration.GRPCAddr)
		if err := grpcServerInstance.Serve(listener); err != nil {
			log.Printf("gRPC serve: %v", err)
		}
	}()

	log.Printf("HTTP gateway listening on %v", configuration.HTTPAddr)
	if err := http.ListenAndServe(configuration.HTTPAddr, httpGateway.Handler()); err != nil {
		log.Fatalf("HTTP: %v", err)
	}
}
