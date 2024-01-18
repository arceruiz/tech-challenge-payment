package main

import (
	"os"
	"os/signal"
	"tech-challenge-payment/internal/channels/rest"
	"tech-challenge-payment/internal/config"
	"tech-challenge-payment/internal/repository"
	"tech-challenge-payment/internal/service"

	"github.com/sirupsen/logrus"
)

var (
	cfg = &config.Cfg
)

func main() {
	config.ParseFromFlags()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		restChannel := rest.NewPaymentChannel(service.NewPaymentService(repository.NewPaymentRepo(repository.NewMongo())))
		if err := rest.New(restChannel).Start(); err != nil {
			logrus.Panic()
		}
	}()

	logrus.WithField("grpc server started on: ", cfg.Server.CustomerPort).Info()
	<-stop
}
