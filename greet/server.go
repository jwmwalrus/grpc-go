package greet

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/grpc-go/greet/greetpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server defines the greet server
type Server struct{}

// GreetWithDeadline implements GreetServiceServer interface
func (s *Server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Printf("GreetWithDeadline function was invoked with request %v", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Error(codes.Canceled, "The client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()

	return &greetpb.GreetWithDeadlineResponse{Result: "Hello " + firstName}, nil
}

// GreetEveryone implements GreetServiceServer interface
func (s *Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error getting request: %v", err)
			return err
		}
		res := &greetpb.GreetEveryoneResponse{Result: "Hello " + req.GetGreeting().GetFirstName()}
		if err := stream.Send(res); err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

// LongGreet implements GreetServiceServer interface
func (s *Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet function was invoked with stream %v", stream)

	list := []string{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &greetpb.LongGreetResponse{Result: strings.Join(list, "\n")}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Printf("Error getting partial request: %v", err)
			return err
		}
		list = append(list, fmt.Sprintf("Hello %s!", req.GetGreeting().GetFirstName()))
	}
}

// GreetManyTimes implementes GreetServiceServer interface
func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{Result: "Hello " + firstName + " number " + strconv.Itoa(i+1)}
		if err := stream.Send(res); err != nil {
			log.Printf("Error sending response #%d of 10: %v", i+1, err)
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// Greet implements GreetServiceServer interface
func (s *Server) Greet(_ context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()

	return &greetpb.GreetResponse{Result: "Hello " + firstName}, nil
}
