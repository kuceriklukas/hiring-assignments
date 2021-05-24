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

// Server is a server implementing the proto image_svc
type Server struct{}

// OptimizeImage performes scaling and greyscaling on the provided image and returns the optimized image
func (s *Server) OptimizeImage(ctx context.Context, req *imagesvc.OptimizeImageRequest) (*imagesvc.OptimizeImageReply, error) {
	fmt.Println("Recieved an image")

	// in the beginning, te new image is the same as the old image (in case no optimizations are done, it will just return the same image)
	newImageBytes := req.GetImage()
	var err error

	// if specified in the request, scale the image
	if req.Scale.GetScale() {
		newImageBytes, err = scaleImage(ctx, newImageBytes, req.Scale.GetTargetWidth(), req.Scale.GetTargetHeight())

		if err != nil {
			fmt.Printf("there was an error: %s", err.Error())
			return &imagesvc.OptimizeImageReply{
				Content: newImageBytes,
			}, err
		}
	}

	// if specified in the request, greyscale the image
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
	// read the bytes and decode them into an image object for further manipulation
	providedImage, _, err := image.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	// resize the image and encode it back to a jpeg image
	newImage := resize.Resize(uint(targetWidth), uint(targetHeight), providedImage, resize.Lanczos3)
	newImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(newImageBuffer, newImage, nil)

	if err != nil {
		return nil, err
	}

	return newImageBuffer.Bytes(), nil
}

func greyScaleImage(ctx context.Context, imageBytes []byte) ([]byte, error) {
	// read the bytes and decode them into an image object for further manipulation
	providedImage, _, err := image.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	// go pixel by pixel on the Y and X axis on the provided image and construct a new "gray image"
	grayImg := image.NewGray(providedImage.Bounds())
	for y := providedImage.Bounds().Min.Y; y < providedImage.Bounds().Max.Y; y++ {
		for x := providedImage.Bounds().Min.X; x < providedImage.Bounds().Max.X; x++ {
			grayImg.Set(x, y, providedImage.At(x, y))
		}
	}

	// read the gray image and encode it into a jpeg image
	newImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(newImageBuffer, grayImg, nil)

	if err != nil {
		return nil, err
	}

	return newImageBuffer.Bytes(), nil
}
