package main

import (
	"context"
	pb "conveyer-service-collector/cmd/collector/measurement"
	"conveyer-service-collector/cmd/collector/model"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// AMQP Channel
var channel *amqp.Channel

const (
	exchange = "conveyer"
)

var grpcServer string = os.Getenv("GRPC_CONN_STR")
var eventQueueConnStr string = os.Getenv("EVENT_QUEUE_CONN_STR")

type service struct {
	pb.UnimplementedMeasurementServiceServer
}

func (s *service) NewMeasurement(ctx context.Context, proto *pb.Measurement) (*pb.OK, error) {

	m := model.NewMeasurement(proto.GetWeight(), proto.GetHumidity(), proto.GetColor())
	json, err := m.MarshallJSON()
	if err != nil {
		log.Fatalf("%s: %s\n", "MarshallJSON", err)
	}

	err = channel.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json,
		})
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to publish a message", err)
	}

	return &pb.OK{
		Success: true,
	}, nil
}

func main() {
	fmt.Println(grpcServer)
	fmt.Println(eventQueueConnStr)
	fmt.Println("- Starting ------")

	conn, err := amqp.Dial(eventQueueConnStr)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	defer channel.Close()

	lis, err := net.Listen("tcp", grpcServer)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
	}

	queue, err := channel.QueueDeclare(
		"clients",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to declare a queue", err)
	}

	err = channel.QueueBind(
		queue.Name,
		"clients",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to `db` queue bind", err)
	}

	queue, err = channel.QueueDeclare(
		"db",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to declare a queue", err)
	}

	err = channel.QueueBind(
		queue.Name,
		"db",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s\n", "Failed to `client` queue bind", err)
	}

	s := grpc.NewServer()
	pb.RegisterMeasurementServiceServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
