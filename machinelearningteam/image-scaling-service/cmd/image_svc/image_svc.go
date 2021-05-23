package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/pkg/image_svc"
	pb "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto/imageoptimizer"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

func main() {
	fmt.Println(fmt.Sprintf("starting the service on the port %s", port))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterImageOptimizerServer(s, &image_svc.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
