package calculator

import (
	"context"
	"errors"
	"io"
	"log"
	"math"

	"github.com/jwmwalrus/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server defines the calculator server
type Server struct{}

// SquareRoot implements CalculatorServiceServer interface
func (s *Server) SquareRoot(_ context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "We don't support complex numbers like %f", number)
	} else if math.IsNaN(number) {
		return nil, errors.New("We don't like NaN")
	}

	return &calculatorpb.SquareRootResponse{Root: math.Sqrt(number)}, nil
}

// FindMaximum implements CalculatorServiceServer interface
func (s *Server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	list := []int32{}
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error getting request: %v", err)
			return err
		}
		list = append(list, req.GetNumber())
		oldMax := max
		max = list[0]
		for i := 1; i < len(list); i++ {
			if list[i] > max {
				max = list[i]
			}
		}
		if oldMax == max {
			continue
		}
		if err := stream.Send(&calculatorpb.FindMaximumResponse{Maximum: max}); err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

// ComputeAverage implements CalculatorServiceServer interface
func (s *Server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := func(list []int32) (acc float64) {
		for _, x := range list {
			acc += float64(x)
		}
		return
	}

	numbers := []int32{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Average: sum(numbers) / float64(len(numbers))})
		}
		if err != nil {
			log.Printf("Error getting request: %v", err)
			return err
		}

		numbers = append(numbers, req.GetNumber())
	}
}

// PrimeFactorization implements CalculatorServiceServer interface
func (s *Server) PrimeFactorization(req *calculatorpb.PrimeFactorizationRequest, stream calculatorpb.CalculatorService_PrimeFactorizationServer) error {
	number := req.GetNumber()

	var k int64 = 2
	for number > 1 {
		if number%k == 0 {
			res := &calculatorpb.PrimeFactorizationResponse{Factor: k}
			if err := stream.Send(res); err != nil {
				log.Printf("Error sending response: %v", err)
				return err
			}
			number = number / k
		} else {
			k++
		}
	}
	return nil
}

// Sum implements CalculatorServiceServer interface
func (s *Server) Sum(_ context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	left := req.GetLeft()
	right := req.GetRight()

	return &calculatorpb.SumResponse{Result: left + right}, nil
}
