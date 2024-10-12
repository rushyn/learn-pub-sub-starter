package pubsub

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error){

	ch, err := conn.Channel()
	if err != nil{
		fmt.Println("Failed to create chanel !!!")
		log.Println(err.Error())
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	var durable = false
	var autoDelete = false
	var exclusive = false

	if simpleQueueType == SimpleQueueDurable {
		durable = true
		autoDelete = false
		exclusive = false
	}else{
		durable = false
		autoDelete = true
		exclusive = true
	}

	args := amqp.Table{}
	args["x-dead-letter-exchange"] = "peril_dlx"

	q, err := ch.QueueDeclare(
								queueName, 			   // queue name
								durable,               // durable
								autoDelete,            // auto-delete
								exclusive,             // exclusive
								false,                 // noWait
								args)

	if err != nil{
		fmt.Println("Failed to create chanel !!!")
		log.Println(err.Error())
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	ch.QueueBind(q.Name, key, exchange, false, nil)
	
	return ch, q, nil

}