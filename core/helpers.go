package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewLogger() (*log.Logger, *os.File, error) {
	var logbuf bytes.Buffer

	logger := log.New(&logbuf, "", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.OpenFile("pixie.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, nil, err
	}

	logger.SetOutput(f)
	logger.Println("Setting up logger")

	return logger, f, nil
}
