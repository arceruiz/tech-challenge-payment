package main

import (
	"tech-challenge-payment/internal/channels/rest"
	"tech-challenge-payment/internal/channels/sqs"
	"tech-challenge-payment/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	config.ParseFromFlags()
	go func() {
		sqs.NewSQS().ReceiveMessage()
	}()

	if err := rest.New().Start(); err != nil {
		logrus.Panic()
	}
}
