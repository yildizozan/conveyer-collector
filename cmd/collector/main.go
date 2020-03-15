package main

import (
	"context"
	pb "conveyer-service-collector/cmd/collector/measurement"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
)

// AMQP Channel
var ch *amqp.Channel

type service struct {
	pb.UnimplementedMeasurementServiceServer
}

func (s *service) NewMeasurement(ctx context.Context, measurement *pb.Measurement) (*pb.OK, error) {
	fmt.Println(measurement)

	err := ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(measurement.String()),
		})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}

	return &pb.OK{
		Success: true,
	}, nil
}

func main() {
	conn, err := amqp.Dial("amqp://orebron.com:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
	}

	s := grpc.NewServer()
	pb.RegisterMeasurementServiceServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
