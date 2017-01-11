package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	c "github.com/dlapiduz/pixie/common"
	"github.com/nlopes/slack"
)

func main() {
	logger, _, _ := c.NewLogger("FAKE")
	ec, recvCh, sendCh := c.ConnectNats()
	defer ec.Close()

	go func() {
		logger.Println("Listening nats")
		for {
			msg := <-recvCh
			logger.Println("New Message")
			fmt.Println(strings.TrimSpace(msg))
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		m := slack.Msg{
			Text: text,
		}
		fmt.Println("Sending Message")
		sendCh <- &m
		fmt.Println("Message sent")

	}
}
