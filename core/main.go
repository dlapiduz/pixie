package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	db.LogMode(false)
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

func RunContainer(client *docker.Client, sendCh chan string, img string, args []string) error {

	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           img,
			AttachStdout:    true,
			AttachStderr:    true,
			NetworkDisabled: false,
			Cmd:             args,
		},
	}

	// opts.Config.Cmd = [""]

	cont, err := client.CreateContainer(opts)
	if err != nil {
		return err
	}

	logger.Println("Data Container ID: " + cont.ID)

	attached := make(chan struct{})

	r, w := io.Pipe()

	attachOptions := docker.AttachToContainerOptions{
		Container:    cont.ID,
		OutputStream: w,
		ErrorStream:  w,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
		Success:      attached,
	}

	logger.Printf("AttachToContainer")
	go client.AttachToContainer(attachOptions)

	go func(reader io.Reader, sendCh chan string) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			sendCh <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
		}
	}(r, sendCh)

	// Wait until
	<-attached
	attached <- struct{}{}

	// start the container
	if err := client.StartContainer(cont.ID, opts.HostConfig); err != nil {
		return err
	}

	logger.Println("Waiting for to exit so we can remove the container\n", cont.ID)
	if _, err := client.WaitContainer(cont.ID); err != nil {
		return err
	}

	removeOpts := docker.RemoveContainerOptions{
		ID: cont.ID,
	}

	logger.Println("Removing container", cont.ID)
	if err := client.RemoveContainer(removeOpts); err != nil {
		return err
	}

	return nil
}
