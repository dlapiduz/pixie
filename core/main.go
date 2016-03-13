package main

import (
	"bytes"
	"fmt"
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

	var f *os.File
	logger, f, _ = c.NewLogger()
	defer f.Close()

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

	fmt.Println("Listening nats")
	for {
		msg := <-recvCh

		if action := RunFilter(db, msg.Text); action.ID > 0 {
			logger.Printf("Running")
			go func() {
				out, err := RunContainer(client, action.Image)
				if err != nil {
					panic(err)
				}

				sendCh <- strings.TrimSpace(out)
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

	db.Where("? ~ trigger", text).First(&action)
	logger.Println(action)

	return action
}

func RunContainer(c *docker.Client, img string) (string, error) {

	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           img,
			AttachStdout:    true,
			AttachStderr:    true,
			NetworkDisabled: false,
		},
		// HostConfig: &docker.HostConfig{
		// 	VolumesFrom: []string{volumesFrom},
		// },
	}

	// opts.Config.Cmd = [""]

	cont, err := c.CreateContainer(opts)
	if err != nil {
		return "", err
	}
	logger.Println("Data Container ID: " + cont.ID)

	attached := make(chan struct{})

	buf := new(bytes.Buffer)

	go func() {
		logger.Printf("AttachToContainer")
		err = c.AttachToContainer(docker.AttachToContainerOptions{
			Container:    cont.ID,
			OutputStream: buf,
			ErrorStream:  os.Stderr,
			Logs:         true,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			Success:      attached,
		})
	}()

	<-attached
	attached <- struct{}{}

	// start the container
	if err := c.StartContainer(cont.ID, opts.HostConfig); err != nil {
		return "", err
	}

	logger.Println("Waiting for to exit so we can remove the container\n", cont.ID)
	if _, err := c.WaitContainer(cont.ID); err != nil {
		return "", err
	}

	removeOpts := docker.RemoveContainerOptions{
		ID: cont.ID,
	}

	logger.Println("Removing container", cont.ID)
	if err := c.RemoveContainer(removeOpts); err != nil {
		return "", err
	}

	return buf.String(), nil
}
