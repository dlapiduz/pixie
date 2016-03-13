package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/nats"
	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_API"))

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		fmt.Println("Error connecting to nats")
		fmt.Println(err)
		return
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	recvCh := make(chan string)
	ec.BindRecvChan("slack.send", recvCh)

	sendCh := make(chan *slack.Msg)
	ec.BindSendChan("slack.receive", sendCh)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var botMeta *slack.UserDetails
	go func() {
		fmt.Println("Listening nats")
		for {
			msg := <-recvCh
			rtm.SendMessage(rtm.NewOutgoingMessage(msg, "C0QJ2EWCT"))
		}
	}()

	fmt.Println("Listening slack")

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				// fmt.Println("Infos:", ev.Info)
				botMeta = ev.Info.User
				_ = botMeta
				// fmt.Println("Connection counter:", ev.ConnectionCount)
				// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C0QJ2EWCT"))
			case *slack.MessageEvent:

				if strings.Contains(ev.Msg.Text, botMeta.ID) {
					sendCh <- &ev.Msg
				}
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
