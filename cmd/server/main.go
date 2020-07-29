package main

import (
	"context"
	pb "github.com/yildizozan/conveyer-collector/pkg/proto/measurement"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type service struct {
	pb.UnimplementedMeasurementServiceServer
}

func (s *service) NewMeasurement(ctx context.Context, measurement *pb.Measurement) (*pb.OK, error) {
	fmt.Println(measurement)

	return &pb.OK{
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMeasurementServiceServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
