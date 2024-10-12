package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)




func main() {
	fmt.Println("Starting Peril client...")
	const conStr = "amqp://guest:guest@localhost:5672/"

	newCon, err := amqp.Dial(conStr)
	if err != nil{
		fmt.Println("amqp Daal Error!!")
		fmt.Printf("%s\n", err.Error())
	}
	defer newCon.Close()

	username, err := gamelogic.ClientWelcome()
	if err !=nil {
		print("You failed to login!!!")
		log.Fatalln(err.Error())
	}

	ch, err := newCon.Channel()

	if err != nil {
		fmt.Println("Something is wrong with newCon.Channel() !!!")
		log.Println(err.Error())
	}


	game := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(newCon, routing.ExchangePerilDirect, routing.PauseKey + "." + username, routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(game))
	if err != nil{
		fmt.Println("SubscribeJSON fail !!!")
		fmt.Println(err.Error())
	}

	err = pubsub.SubscribeJSON(newCon, routing.ExchangePerilTopic, routing.ArmyMovesPrefix + "." + username, routing.ArmyMovesPrefix + ".*", pubsub.SimpleQueueTransient, handlerMove(game, ch))
	if err != nil{
		fmt.Println("SubscribeJSON fail !!!")
		fmt.Println(err.Error())
	}

	err = pubsub.SubscribeJSON(newCon, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix + ".*", pubsub.SimpleQueueDurable, handlerWar(game))
	if err != nil{
		fmt.Println("SubscribeJSON fail !!!")
		fmt.Println(err.Error())
	}



	validLocations := []string{
		"americas",
		"europe",
		"africa",
		"asia",
		"antarctica",
		"australia",
	}

	validUnits := []string{
		"infantry",
		"cavalry",
		"artillery",
	}
	

	for{

		input := gamelogic.GetInput()

		if input[0] == "spawn"{

			validLocation 	:= slices.Contains(validLocations, input[1])
			validUnit 		:= slices.Contains(validUnits, input[2])
			
			if !validLocation{
				log.Printf("Invalid locatuib %s \n", input[1])
			}
			if !validUnit{
				log.Printf("Invalid unit type %s \n", input[2])
			}

			if validLocation && validUnit {
				err := game.CommandSpawn(input)
				if err != nil {
					log.Println("failed to Spawn")
					log.Println(err.Error())
				}

			}
			continue
		}


		if input[0] == "move"{
			move, err := game.CommandMove(input)
			if err != nil {
				log.Println("failed to move")
				log.Println(err.Error())
			}else{
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix + "." + username, move)
				if err != nil{
					log.Println("Unable to publish move.")
				}else{
					log.Println("!!! Move was published !!!")
				}
			}
			continue
		}

		if input[0] == "status"{
			game.CommandStatus()
			continue
		}

		if input[0] == "spam"{
			fmt.Println("Spamming not allowed yet!")
			continue
		}

		if input[0] == "quit"{
			gamelogic.PrintQuit()
			continue
		}

		log.Printf("%s is an invalid command \n", input[0])
		
	}

}
