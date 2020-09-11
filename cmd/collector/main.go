package main

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/yildizozan/conveyor-collector/cmd/model"
	pb "github.com/yildizozan/conveyor-collector/pkg/proto/conveyor"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// AMQP Channel
var channel *amqp.Channel

const (
	exchange = "conveyor"
)

var grpcServer = os.Getenv("GRPC_CONN_STR")
var eventQueueConnStr = os.Getenv("EVENT_QUEUE_CONN_STR")

type service struct {
	pb.UnimplementedConveyorServiceServer
}

func (s *service) NewPosition(ctx context.Context, proto *pb.Position) (*pb.OK, error) {

	m := model.NewPosition(proto.GetX(), proto.GetY(), proto.GetZ())
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

func (s *service) NewEnginesState(ctx context.Context, proto *pb.Engines) (*pb.OK, error) {

	m := model.NewEngines(proto.GetGreen(), proto.GetRed(), proto.GetBlack(), proto.GetOrange())
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

	// Exchange
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

	// New queue declare
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

	s := grpc.NewServer()
	pb.RegisterConveyorServiceServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
