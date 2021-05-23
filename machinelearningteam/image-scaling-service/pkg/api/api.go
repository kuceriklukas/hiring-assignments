package api

import (
	"context"
	"fmt"
	"log"

	api "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto"
	image_svc "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto/imageoptimizer"
	"google.golang.org/grpc"
)

// Server is a server implementing the proto API
type Server struct{}

const (
	host = "localhost:50052"
)

// ScaleImage echoes the image provides in the request
func (s *Server) ScaleImage(ctx context.Context, req *api.ScaleImageRequest) (*api.ScaleImageReply, error) {
	// Echo
	fmt.Println("Recieved image...")

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := image_svc.NewImageOptimizerClient(conn)

	image := req.Image.GetContent()

	resp, err := client.OptimizeImage(ctx, &image_svc.OptimizeImageRequest{
		Image:     image,
		Scale:     true,
		Greyscale: true,
	})

	if err != nil {
		fmt.Printf("Error getting file %+v\n", resp)
	}

	return &api.ScaleImageReply{
		Content: resp.GetContent(),
	}, nil
}
