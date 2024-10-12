package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	//log.Println("Publishing Message to Ch")
	Jval, err := json.Marshal(val)
	if err != nil{
		fmt.Println("fail to marshal $val")
		fmt.Println(err.Error())
		return err
	}
	
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType:     "application/json",
		Body:            Jval,
	})
	if err != nil{
		fmt.Println("fail to PublishWithContext")
		fmt.Println(err.Error())
		return err
	}

	return nil
}