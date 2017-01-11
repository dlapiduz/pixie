package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

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

	logger.Println("Args: ", args)

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
