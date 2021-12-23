package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/jwmwalrus/grpc-go/blog"
	"github.com/jwmwalrus/grpc-go/blog/blogpb"
	"github.com/jwmwalrus/grpc-go/calculator"
	"github.com/jwmwalrus/grpc-go/calculator/calculatorpb"
	"github.com/jwmwalrus/grpc-go/greet"
	"github.com/jwmwalrus/grpc-go/greet/greetpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	useTLS   = false
	certFile = "ssl/server.crt"
	keyFile  = "ssl/server.pem"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Running server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Print("Listening on port 50051")

	// serveGreeting(lis)
	// serveCalculator(lis)
	serveBlog(lis)
}

func serveBlog(lis net.Listener) {
	log.Print("Serving Blog...")

	log.Print("Connecting to database...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	srv := &blog.Server{}
	srv.Coll = client.Database("mydb").Collection("blog")

	opts := getOpts()
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, srv)
	reflection.Register(s)

	go func() {
		log.Print("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Print("Stopping the server")
	s.Stop()
	log.Print("Closing the listener")
	lis.Close()
	log.Print("Disconnecting mongodb")
	client.Disconnect(context.TODO())
	log.Print("End of program")
}

func serveCalculator(lis net.Listener) {
	log.Print("Serving Calculator...")

	opts := getOpts()
	s := grpc.NewServer(opts...)
	calculatorpb.RegisterCalculatorServiceServer(s, &calculator.Server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func serveGreeting(lis net.Listener) {
	log.Print("Serving Greet...")

	opts := getOpts()
	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &greet.Server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getOpts() (opts []grpc.ServerOption) {
	if useTLS {
		log.Print("Using TLS")
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Error obtaining credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	return
}
