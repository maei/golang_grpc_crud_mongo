package main

import (
	"context"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/app"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/client"
	"github.com/maei/golang_grpc_crud_mongo/grpc_server/src/server"
	"github.com/maei/shared_utils_go/logger"
	"sync"
)

func main() {
	logger.Info("gRPC-Server: Starting MongoDB connection")
	app.StartApplication()
	logger.Info("gRPC-Server: Starting gRPC-Server")
	var wg sync.WaitGroup

	wg.Add(1)

	go server.StartGRPCServer(&wg)

	// End gRPC-Server gracefully
	wg.Wait()
	logger.Info("gRPC-Server: Disconnecting MongoDB Client")
	client.MongoClient.Disconnect(context.Background())

}
