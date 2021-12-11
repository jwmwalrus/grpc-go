package greet

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/jwmwalrus/grpc-go/greet/greetpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoDeadline implerments a unary invocation with deadline
func DoDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	log.Println("Starting to do a unary client RPC for GreetWithDeadline...")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "M",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.DeadlineExceeded {
			log.Printf("Call is taking too much time: %s", s.Message())
			return
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v", s.Message())
		}
	}

	log.Printf("Response from server: %s", res.GetResult())
}

// DoBidirectionalStreaming implerments a bidirectional streaming  invocation
func DoBidirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a streaming client RPC for GreetEveryone...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC:%v", err)
	}

	names := []string{"Dave Dee", "Dozy", "Beaky", "Mick", "Tich"}

	waitChannel := make(chan struct{})

	go func() {
		for _, first := range names {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: first,
					LastName:  "M",
				},
			}
			log.Printf("Sending request: %v", req)
			if err := stream.Send(req); err != nil {
				log.Fatalf("Error while calling GreetEveryone RPC:%v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		defer close(waitChannel)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while getting response for GreetEveryone RPC:%v", err)
			}
			log.Printf("Response from GreetEveryone: \n%v", res.GetResult())
		}
	}()

	<-waitChannel
}

// DoClientStreaming implements a streaming client invocation
func DoClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a streaming client RPC for LongGreet...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC:%v", err)
	}
	names := []string{"Dave Dee", "Dozy", "Beaky", "Mick", "Tich"}
	for _, first := range names {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: first,
				LastName:  "M",
			},
		}
		log.Printf("Sending request: %v", req)
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while calling GreetManyTimes RPC:%v", err)
		}
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Printf("Error while getting response from LongGreet RPC:%v", err)
	}
	log.Printf("Response from GreetManyTimes: \n%v", res.GetResult())
}

// DoServerStreaming implements a streaming server invocation
func DoServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a streaming server RPC for GreetManyTimes...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "M",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC:%v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while calling GreetManyTimes RPC:%v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", res.GetResult())
	}
}

// DoUnary implements a unary client invocation
func DoUnary(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a unary RPC for Greet...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "M",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC:%v", err)
	}
	log.Printf("Response from Greet: %v", res.GetResult())
}
