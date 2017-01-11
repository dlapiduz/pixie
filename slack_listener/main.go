package main

import (
	"os"
	"strings"

	c "github.com/dlapiduz/pixie/common"
	"github.com/nlopes/slack"
)

func main() {
	logger, f, _ := c.NewLogger("SLACK")
	defer f.Close()

	api := slack.New(os.Getenv("SLACK_API"))

	ec, recvCh, sendCh := c.ConnectNats()
	defer ec.Close()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var botMeta *slack.UserDetails
	go func() {
		logger.Println("Listening nats")
		for {
			msg := <-recvCh
			logger.Println("Sending Message")
			rtm.SendMessage(rtm.NewOutgoingMessage(strings.TrimSpace(msg), "C0QJ2EWCT"))
			logger.Println("Message Sent")
		}
	}()

	logger.Println("Listening slack")

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
				logger.Println("Received message")
				if strings.Contains(ev.Msg.Text, botMeta.ID) {
					logger.Println("Message matched")
					sendCh <- &ev.Msg
					logger.Println("Message sent")
				}
			case *slack.PresenceChangeEvent:
				// fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				// fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				logger.Println("Error: ", ev.Error())

			case *slack.InvalidAuthEvent:
				logger.Println("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
