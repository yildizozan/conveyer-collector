package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// AMQP Channel
var ch *amqp.Channel

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

	rand.Seed(time.Now().UTC().UnixNano())
	for x := range time.Tick(3 * time.Second) {
		fmt.Printf("%s\n", x)

		body := bodyFrom(os.Args)
		err = ch.Publish(
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalf("%s: %s", "Failed to publish a message", err)
		}
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = fmt.Sprintf("Humidity: %f", rand.Float64())
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
