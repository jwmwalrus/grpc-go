package main

import (
	"fmt"
	"log"

	"github.com/jwmwalrus/grpc-go/blog"
	"github.com/jwmwalrus/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	useTLS   = false
	certFile = "ssl/ca.crt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Running client")

	opts := grpc.WithInsecure()
	if useTLS {
		log.Print("Using TLS")
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatal("Eror obtaining credentials", err)
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	// c := greetpb.NewGreetServiceClient(cc)
	// greet.DoUnary(c)
	// greet.DoServerStreaming(c)
	// greet.DoClientStreaming(c)
	// greet.DoBidirectionalStreaming(c)
	// greet.DoDeadline(c, 1*time.Second)
	// greet.DoDeadline(c, 5*time.Second)

	// c := calculatorpb.NewCalculatorServiceClient(cc)
	// calculator.DoUnary(c)
	// calculator.DoStreamingServer(c)
	// calculator.DoStreamingClient(c)
	// calculator.DoBidirectionalStreaming(c)
	// calculator.DoErrorUnary(c)

	c := blogpb.NewBlogServiceClient(cc)
	// blog.DoCreateBlog(c)
	// blog.DoReadBlog(c)
	// blog.DoUpdateBlog(c)
	// blog.DoDeleteBlog(c)
	blog.DoListBlog(c)
}
