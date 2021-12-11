package calculator

import (
	"context"
	"io"
	"log"
	"math"

	"github.com/jwmwalrus/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoErrorUnary implerments a unary  invocation with error handling
func DoErrorUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a unary RPC for SquareRoot...")

	numbers := []float64{14.14, 16.0, -64.0, 225., math.NaN(), 0.0}
	for _, n := range numbers {
		log.Printf("For number: %g", n)
		res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
		if err != nil {
			s := status.Convert(err)
			if s.Code() == codes.InvalidArgument {
				log.Printf("\tError computing square root: %v", s.Message())
			} else {
				log.Fatalf("\tError obtaining response: %v", s.Message())
			}
			continue
		}
		log.Printf("\tRoot is: %g", res.GetRoot())
	}
}

// DoBidirectionalStreaming implerments a bidirectional streaming  invocation
func DoBidirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a streaming client RPC for FindMaximum...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error invoking FindMaximum: %v", err)
	}

	numbers := []int32{1, 5, 3, 6, 2, 20}

	waitChannel := make(chan struct{})

	go func() {
		for _, n := range numbers {
			req := &calculatorpb.FindMaximumRequest{Number: n}
			log.Printf("Sending request: %v", req)
			if err := stream.Send(req); err != nil {
				log.Fatal("Error sending request:", err)
			}
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
				log.Fatalf("Error obtaining response: %v", err)
			}
			log.Printf("Maximum is now: %v", res.GetMaximum())
		}
	}()

	<-waitChannel
}

// DoStreamingClient implements a streaming server client invocation
func DoStreamingClient(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a streaming client RPC for ComputeAverage...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage RPC:%v", err)
	}
	for i := 1; i <= 4; i++ {
		req := &calculatorpb.ComputeAverageRequest{Number: int32(i)}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while senting partial RPC request:%v", err)
		}

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while obtaining RPC responset:%v", err)
	}
	log.Printf("Average is: %g", res.GetAverage())

}

// DoStreamingServer implements a streaming server client invocation
func DoStreamingServer(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a streaming server RPC for PrimeFactorization...")

	req := &calculatorpb.PrimeFactorizationRequest{Number: 2525}
	log.Printf("Number is: %d", req.GetNumber())

	stream, err := c.PrimeFactorization(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeFactorization RPC:%v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while obtaining factor from RPC call:%v", err)
		}
		log.Printf("\tOne factor is: %v", res.GetFactor())
	}
}

// DoUnary implements a unary client invocation
func DoUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a unary RPC for Sum...")

	req := &calculatorpb.SumRequest{Left: 1, Right: 2}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC:%v", err)
	}
	log.Printf("Response from Sum: %v", res.GetResult())
}
