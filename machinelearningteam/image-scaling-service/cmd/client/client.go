package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/kuceriklukas/hiring-assigments/machinelearningteam/image-scaling-service/proto"
	"google.golang.org/grpc"
)

const (
	host = "localhost:50051"
)

func main() {
	fmt.Printf("starting the client and connecting to server at %s\n", host)

	// read and parse the config flags from when the app was run
	var url string
	flag.StringVar(&url, "url", "", "Provide the url for downloading the image from instead of using the defualt test.jpg")
	flag.Parse()

	// connect to the GRPC server
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewImageScalerClient(conn)

	var resp *pb.ScaleImageReply
	ctx := context.Background()

	// if the url was specified, take that into the request for the api
	if url != "" {
		resp, err = client.ScaleImage(ctx, &pb.ScaleImageRequest{
			Image: &pb.Image{
				Source: &pb.ImageSource{
					HttpUri: url, // e.g. "https://place-puppy.com/800x650"
				},
			},
		})

		if err != nil {
			log.Fatal("Couldn't fetch image based on given url")
		}
	} else {
		// if the url wasn't specified, take the test file into the request for the api
		image, err := ioutil.ReadFile("test.jpg")
		if err != nil {
			log.Fatal("Couldn't read input image")
		}

		resp, err = client.ScaleImage(ctx, &pb.ScaleImageRequest{
			Image: &pb.Image{
				Content: image,
			},
		})

		if err != nil {
			fmt.Println("Couldn't fetch image based on the test image")
		}
	}

	// write out the file
	ioutil.WriteFile("out.jpg", resp.GetContent(), 0644)
}
