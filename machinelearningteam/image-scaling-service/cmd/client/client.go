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

	var url string
	flag.StringVar(&url, "url", "", "Provide the url for downloading the image from instead of using the defualt test.jpg")
	flag.Parse()

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewImageScalerClient(conn)

	var resp *pb.ScaleImageReply
	ctx := context.Background()

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

	ioutil.WriteFile("out.jpg", resp.GetContent(), 0644)
}
