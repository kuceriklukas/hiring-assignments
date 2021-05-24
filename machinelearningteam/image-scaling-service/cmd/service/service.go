package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/pkg/api"
	pb "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	health "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/pkg/health/v1"
	api_health "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	port = ":50051"
)

// starts the Prometheus stats endpoint server
func startPromHTTPServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println("prometheus err", port)
	}
}

func main() {
	fmt.Printf("starting the server on the port %s\n", port)
	var imageTargetWidth int
	var imageTargetHeight int
	var doGreyscale bool

	flag.IntVar(&imageTargetWidth, "w", 1024, "The desired image width after scaling, default being 1024")
	flag.IntVar(&imageTargetHeight, "h", 768, "The desired image height after scaling, default being 768")
	flag.BoolVar(&doGreyscale, "gs", true, "Whether the api should greyscale the image as well or not, default being true")
	flag.Parse()

	fmt.Printf("starting the api with follwoing config: target width = %d, target height = %d, perform greyscaling = %t\n", imageTargetWidth, imageTargetHeight, doGreyscale)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go startPromHTTPServer("5001")

	s := grpc.NewServer()

	// Register: Health
	healthServ := health.NewHealthCheckService()
	api_health.RegisterHealthServer(s, healthServ)

	apiServer := &api.Server{
		TargetWidth:  imageTargetWidth,
		TargetHeight: imageTargetHeight,
		DoGreyscale:  doGreyscale,
	}

	pb.RegisterImageScalerServer(s, apiServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
