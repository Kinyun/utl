package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"google.golang.org/grpc"
)

func StartGrpc(grpcServer *grpc.Server, lis net.Listener) {
	// Start server
	fmt.Println("Serving gRPC on tcp://" + lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("StartGrpc - failed to serve: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Stopping gRPC Server")
}
