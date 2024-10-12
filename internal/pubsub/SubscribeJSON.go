package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)


func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // // an enum to represent "durable" or "transient"
	handler func(T) (Acktype),
) error {
	c, q, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ch, err := c.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	go func() {
		defer c.Close()
		for msg := range ch{
			var gen T
			json.Unmarshal(msg.Body, &gen)
			responce := handler(gen)

			switch responce {
				case Ack:
					msg.Ack(false)
					log.Println("Ack")
				case NackRequeue:
					msg.Nack(false, true)
					log.Println("Nack Requeue")
				case NackDiscard:
					msg.Nack(false, false)
					log.Println("Nack Discard")
			}

		}
	}()

	return nil
}