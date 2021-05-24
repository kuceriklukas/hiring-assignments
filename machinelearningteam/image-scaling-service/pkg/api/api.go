package api

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"net/http"

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
	host          = "localhost:50052"
	scalingMargin = 5
)

// ScaleImage calculates whether the image needs to be scaled or not and forwards this to image_svc and then returns the resulting image
func (s *Server) ScaleImage(ctx context.Context, req *api.ScaleImageRequest) (*api.ScaleImageReply, error) {
	fmt.Println("Received an image")

	// make connection to the image_svc via grpc
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := image_svc.NewImageOptimizerClient(conn)

	// get the images bytes, either straight from the request or from the url
	var imageBytes []byte
	if req.Image.Source.GetHttpUri() != "" {
		imageBytes, err = downloadImageFromUri(req.Image.Source.GetHttpUri())
		if err != nil {
			fmt.Printf("Error downloading file:%s", err.Error())
		}
	} else {
		imageBytes = req.Image.GetContent()
	}

	// create the image object for later we need to extract some information from it
	providedImage, _, err := image.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		fmt.Printf("Error decoding image %+v\n", err.Error())
	}

	// calculate if we need to scale the image and what the height and width should be
	scalingOptions := calculateScalingOptions(providedImage, s.TargetWidth, s.TargetHeight)

	// create the response object
	resp, err := client.OptimizeImage(ctx, &image_svc.OptimizeImageRequest{
		Image:     imageBytes,
		Scale:     scalingOptions,
		Greyscale: s.DoGreyscale,
	})

	if err != nil {
		fmt.Printf("Error getting file %+v\n", resp)
	}

	return &api.ScaleImageReply{
		Content: resp.GetContent(),
	}, nil
}

func calculateScalingOptions(image image.Image, targetWidth int, targetHeight int) *image_svc.SizingOptions {
	actualTargetWidth := image.Bounds().Max.X
	actualTargetHeight := image.Bounds().Max.Y
	scaleImage := false

	// calculate what the max and min widths and heights are (depending on the scalingMargin in percentage)
	// since division might end up as a float, we need to careful with the types here, that's why so many conversions
	widthMax := int(float32(targetWidth) / float32(100) * (float32(100 + scalingMargin)))
	widthMin := int(float32(targetWidth) / float32(100) * (float32(100 - scalingMargin)))
	heightMax := int(float32(targetHeight) / float32(100) * (float32(100 + scalingMargin)))
	heightMin := int(float32(targetHeight) / float32(100) * (float32(100 - scalingMargin)))

	// check if the image width is bigger or smaller than the desired maxWidth or minWidth
	if image.Bounds().Max.X < widthMin || image.Bounds().Max.X > widthMax {
		actualTargetWidth = targetWidth
		scaleImage = true

		// if the image width has been scaled, we need to scale the height of the image as well, so we don't destroy the image's ratio
		if scaleImage {
			scalingRatio := float32(targetWidth) / float32(image.Bounds().Max.X)
			actualTargetHeight = int(float32(image.Bounds().Max.Y) * scalingRatio)
		}
	} else if image.Bounds().Max.Y < heightMin || image.Bounds().Max.Y > heightMax {
		// check if the image height is bigger or smaller than the desired maxHeight or minHeight
		actualTargetHeight = targetHeight
		scaleImage = true

		// if the image width has been scaled, we need to scale the width of the image as well, so we don't destroy the image's ratio
		if scaleImage {
			scalingRatio := float32(targetHeight) / float32(image.Bounds().Max.Y)
			actualTargetWidth = int(float32(image.Bounds().Max.X) * scalingRatio)
		}
	}

	return &image_svc.SizingOptions{
		Scale:        scaleImage,
		TargetWidth:  int32(actualTargetWidth),
		TargetHeight: int32(actualTargetHeight),
	}
}

func downloadImageFromUri(uri string) ([]byte, error) {
	// fetch the image from the given url
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	// check if the request to fetch the image was successful
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("couldn't download image, received status code %d", resp.StatusCode)
	}

	// read the response's body into a byte array
	defer resp.Body.Close()
	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
