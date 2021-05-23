package image_svc

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"

	imagesvc "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto/imageoptimizer"
	"github.com/nfnt/resize"
)

// Server is a server implementing the proto API
type Server struct{}

// ScaleImage echoes the image provides in the request
func (s *Server) OptimizeImage(ctx context.Context, req *imagesvc.OptimizeImageRequest) (*imagesvc.OptimizeImageReply, error) {
	// Echo
	fmt.Println("Recieved an image")

	providedImage, _, err := image.Decode(bytes.NewReader(req.GetImage()))

	if err != nil {
		fmt.Printf("there was an error: %s", err.Error())
		return &imagesvc.OptimizeImageReply{
			Content: nil,
		}, err
	}

	newImage := resize.Resize(100, 200, providedImage, resize.Lanczos3)

	newImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(newImageBuffer, newImage, nil)

	if err != nil {
		fmt.Printf("there was an error: %s", err.Error())
		return &imagesvc.OptimizeImageReply{
			Content: nil,
		}, err
	}

	return &imagesvc.OptimizeImageReply{
		Content: newImageBuffer.Bytes(),
	}, nil
}
