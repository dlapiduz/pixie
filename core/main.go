package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	c "github.com/dlapiduz/pixie/common"
	"github.com/fsouza/go-dockerclient"
)

type Action struct {
	Trigger string
	Image   string
}

var Actions = []Action{{"hello", "dlapiduz/pixie-hello"}}

var logger *log.Logger

func main() {
	db, _ := c.LoadDB()
	db.DB().Ping()

	client, _ := docker.NewClientFromEnv()

	var f *os.File
	logger, f, _ = NewLogger()
	defer f.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		if img := RunFilter(text); img != "" {
			logger.Printf("Running")
			go func() {
				out, err := RunContainer(client, img)
				if err != nil {
					panic(err)
				}

				fmt.Println(out)
			}()

		}

	}
}

func RunFilter(text string) string {
	for _, a := range Actions {
		rp, err := regexp.Compile(a.Trigger)
		if err != nil {
			panic("compile error")
		}

		if rp.MatchString(text) {
			return a.Image
		}
	}
	return ""
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
