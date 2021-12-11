package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jwmwalrus/grpc-go/calculator"
	"github.com/jwmwalrus/grpc-go/calculator/calculatorpb"
	"github.com/jwmwalrus/grpc-go/greet"
	"github.com/jwmwalrus/grpc-go/greet/greetpb"
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
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Print("Listening on port 50051")

	// serveGreeting(lis)
	serveCalculator(lis)
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
