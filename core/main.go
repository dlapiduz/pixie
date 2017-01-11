package main

import (
	"fmt"
	"log"
	"strings"

	c "github.com/dlapiduz/pixie/common"
	"github.com/fsouza/go-dockerclient"
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

type Action struct {
	Trigger string
	Image   string
}

var Actions = []Action{{"hello", "dlapiduz/pixie-hello"}}

var logger *log.Logger

func main() {
	db, _ := c.LoadDB()
	db.LogMode(true)
	db.DB().Ping()

	client, _ := docker.NewClientFromEnv()

	logger, _, _ = c.NewLogger("CORE")
	// defer f.Close()

	// Create nats connection
	ec := c.CreateNatsConn(logger)
	defer ec.Close()

	if ec == nil {
		return
	}

	recvCh := make(chan *slack.Msg)
	ec.BindRecvChan("slack.receive", recvCh)

	sendCh := make(chan string)
	ec.BindSendChan("slack.send", sendCh)

	logger.Println("Listening nats")
	for {
		msg := <-recvCh
		logger.Println("Received message")
		if action := RunFilter(db, msg.Text); action.ID > 0 {
			logger.Printf("Running")
			go func() {
				args := strings.Split(strings.Trim(action.Match, "{}"), ",")
				err := RunContainer(client, sendCh, action.Image, args)
				if err != nil {
					logger.Println("Error running container")
					logger.Println(err)
				}

				logger.Printf("Sent")
			}()
		}

	}
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("Enter text: ")
	// 	text, _ := reader.ReadString('\n')
	// 	text = strings.TrimSpace(text)

	// 	}

	// }
}

func RunFilter(db *gorm.DB, text string) c.Action {
	var action c.Action
	if text == "" {
		return action
	}

	query := fmt.Sprintf("id, trigger, image, regexp_matches('%s', trigger) as match", text)

	db.Select(query).Where("? ~ trigger", text).First(&action)
	logger.Println(action)

	return action
}
