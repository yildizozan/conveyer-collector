// Package main implements a sensor for Greeter service.
package main

import (
	"context"
	pb "github.com/yildizozan/conveyer-collector/pkg/proto/measurement"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMeasurementServiceClient(conn)

	// Important for setTimeout
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//	defer cancel()

	// Starting send data
	rand.Seed(time.Now().UTC().UnixNano())
	for x := range time.Tick(3 * time.Second) {

		measurement := &pb.Measurement{
			Weight:   rand.Float32(),
			Humidity: rand.Float32(),
			Color:    rand.Float32(),
		}

		reply, err := client.NewMeasurement(context.Background(), measurement)
		if err != nil {
			log.Fatalf("%v.NewMeasurement(_) = _, %v", client, err)
		}

		fmt.Printf("%v %v\n", x, reply)
	}
}
