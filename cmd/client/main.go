package main

import (
	"tech-challenge-payment/internal/channels/rest"
	"tech-challenge-payment/internal/channels/sqs"
	"tech-challenge-payment/internal/config"

	"github.com/rs/zerolog/log"
)

func main() {
	config.ParseFromFlags()

	log.Info().Any("config", config.Get()).Msg("configuration file")
	go func() {
		sqs.NewSQS().ReceiveMessage()
	}()

	if err := rest.New().Start(); err != nil {
		log.Fatal().Err(err).Msg("an error occurred in rest channel")
	}
}
