package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	gamelogic.PrintServerHelp()

	const conStr = "amqp://guest:guest@localhost:5672/"

	newCon, err := amqp.Dial(conStr)
	if err != nil{
		fmt.Println("amqp Daal Error!!")
		fmt.Printf("%s\n", err.Error())
	}
	defer newCon.Close()
	fmt.Println("Peril server Started")

	ch, err := newCon.Channel()
	if err != nil{
		fmt.Println("channel Error!!")
		fmt.Printf("%s\n", err.Error())
	}

	for {
		input := gamelogic.GetInput()

		if input[0] == "pause"{

			log.Println("pause")

			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				fmt.Println("PublishJason Error")
				fmt.Println(err.Error())
			}
			continue
		}

		if input[0] == "resume"{

			log.Println("resume")
			
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				fmt.Println("PublishJason Error")
				fmt.Println(err.Error())
			}

			continue

		}

		if input[0] == "quit"{

			log.Println("quit")
			break
		}


		log.Printf("%s is an invalid command \n", input[0])
		



	}

}
