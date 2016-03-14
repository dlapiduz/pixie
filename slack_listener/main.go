package main

import (
	"os"
	"strings"

	c "github.com/dlapiduz/pixie/common"
	"github.com/nats-io/nats"
	"github.com/nlopes/slack"
)

func main() {
	logger, f, _ := c.NewLogger("SLACK")
	defer f.Close()

	api := slack.New(os.Getenv("SLACK_API"))

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		logger.Println("Error connecting to nats")
		logger.Println(err)
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
