package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const conStr = "amqp://guest:guest@localhost:5672/"

	newCon, err := amqp.Dial(conStr)
	if err != nil{
		fmt.Println("amqp Daal Error!!")
		fmt.Printf("%s\n", err.Error())
	}
	defer newCon.Close()
	fmt.Println("Peril server Started")
	fmt.Println("----------------------------------")
	fmt.Println("Wait for a signal (Ctrl+C) to exit")
	fmt.Println("----------------------------------")

	ch, err := newCon.Channel()
	if err != nil{
		fmt.Println("channel Error!!")
		fmt.Printf("%s\n", err.Error())
	}


	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		fmt.Println("PublishJason Error")
		fmt.Println(err.Error())
	}

}
