package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_API"))
	// api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var botMeta *slack.UserDetails

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				// fmt.Println("Infos:", ev.Info)
				botMeta = ev.Info.User
				_ = botMeta
				// fmt.Println("Connection counter:", ev.ConnectionCount)
				// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C0QJ2EWCT"))
			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				// if ev.Members = botMeta.ID {
				// 	rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C0QJ2EWCT"))
				// }
			case *slack.PresenceChangeEvent:
				// fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				// fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
