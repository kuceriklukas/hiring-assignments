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
	fmt.Println("Recieved an image")

	newImageBytes := req.GetImage()
	var err error

	if req.Scale.GetScale() {
		newImageBytes, err = scaleImage(ctx, newImageBytes, req.Scale.GetTargetWidth(), req.Scale.GetTargetHeight())

		if err != nil {
			fmt.Printf("there was an error: %s", err.Error())
			return &imagesvc.OptimizeImageReply{
				Content: newImageBytes,
			}, err
		}
	}

	if req.GetGreyscale() {
		newImageBytes, err = greyScaleImage(ctx, newImageBytes)

		if err != nil {
			fmt.Printf("there was an error: %s", err.Error())
			return &imagesvc.OptimizeImageReply{
				Content: newImageBytes,
			}, err
		}
	}

	return &imagesvc.OptimizeImageReply{
		Content: newImageBytes,
	}, nil
}

func scaleImage(ctx context.Context, imageBytes []byte, targetWidth int32, targetHeight int32) ([]byte, error) {
	providedImage, _, err := image.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	newImage := resize.Resize(uint(targetWidth), uint(targetHeight), providedImage, resize.Lanczos3)
	newImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(newImageBuffer, newImage, nil)

	if err != nil {
		return nil, err
	}

	return newImageBuffer.Bytes(), nil
}

func greyScaleImage(ctx context.Context, imageBytes []byte) ([]byte, error) {
	providedImage, _, err := image.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	grayImg := image.NewGray(providedImage.Bounds())
	for y := providedImage.Bounds().Min.Y; y < providedImage.Bounds().Max.Y; y++ {
		for x := providedImage.Bounds().Min.X; x < providedImage.Bounds().Max.X; x++ {
			grayImg.Set(x, y, providedImage.At(x, y))
		}
	}

	newImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(newImageBuffer, grayImg, nil)

	if err != nil {
		return nil, err
	}

	return newImageBuffer.Bytes(), nil
}
