package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + "5672"
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return channel

}

func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {
	err := ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Menssagem enviada")

}

func StartConsuming(ch *amqp.Channel, in chan []byte) {
	q, err := ch.QueueDeclare(
		"checkout_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "checkout", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for m := range msgs {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
