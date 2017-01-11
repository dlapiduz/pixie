package common

import (
	"os"

	"github.com/nats-io/nats"
	"github.com/nlopes/slack"
)

func ConnectNats() (*nats.EncodedConn, chan string, chan *slack.Msg) {
	logger, _, _ := NewLogger("NATS")
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		logger.Println("Error connecting to nats")
		logger.Println(err)
		return nil, nil, nil
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	recvCh := make(chan string)
	ec.BindRecvChan("slack.send", recvCh)

	sendCh := make(chan *slack.Msg)
	ec.BindSendChan("slack.receive", sendCh)

	return ec, recvCh, sendCh
}
