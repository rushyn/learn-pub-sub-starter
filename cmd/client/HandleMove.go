package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)





func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove)(pubsub.Acktype){

	return func(am gamelogic.ArmyMove) (pubsub.Acktype) {
		defer fmt.Print("> ")
		mo := gs.HandleMove(am)

		switch mo {
			case gamelogic.MoveOutcomeSamePlayer:
				return pubsub.NackDiscard
			case gamelogic.MoveOutComeSafe:
				return pubsub.Ack
			case gamelogic.MoveOutcomeMakeWar:
				msg := gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.GetPlayerSnap(),
				}
				err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix + "." + gs.GetUsername(), msg)
				if err != nil{
					log.Println(err.Error())
					return pubsub.NackRequeue
				}
				return pubsub.Ack
			default:
				return pubsub.NackDiscard
		}
	}

}