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
type Server struct {
	TargetWidth  int
	TargetHeight int
	DoGreyscale  bool
}

const (
	host = "localhost:50052"
)

// ScaleImage echoes the image provides in the request
func (s *Server) ScaleImage(ctx context.Context, req *api.ScaleImageRequest) (*api.ScaleImageReply, error) {
	fmt.Println("Received image...")

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := image_svc.NewImageOptimizerClient(conn)

	image := req.Image.GetContent()

	resp, err := client.OptimizeImage(ctx, &image_svc.OptimizeImageRequest{
		Image: image,
		Scale: &image_svc.SizingOptions{
			Scale:        true,
			TargetWidth:  int32(s.TargetWidth),
			TargetHeight: int32(s.TargetHeight),
		},
		Greyscale: s.DoGreyscale,
	})

	if err != nil {
		fmt.Printf("Error getting file %+v\n", resp)
	}

	return &api.ScaleImageReply{
		Content: resp.GetContent(),
	}, nil
}
