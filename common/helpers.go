package common

import (
	"bytes"
	"log"
	"math/rand"
	"os"

	"github.com/nats-io/nats"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewLogger(prefix string) (*log.Logger, *os.File, error) {
	var logbuf bytes.Buffer

	logger := log.New(&logbuf, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	f, err := os.OpenFile("../pixie.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, nil, err
	}

	logger.SetOutput(f)
	logger.Println("Setting up logger")

	logger.SetPrefix(prefix + ": ")

	return logger, f, nil
}

func CreateNatsConn(logger *log.Logger) *nats.EncodedConn {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		logger.Println("Error connecting to nats")
		logger.Println(err)
		return nil
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	return ec
}
