package main

import (
	"net/http"
	"tech-challenge-payment/internal/channels/rest"
	"tech-challenge-payment/internal/config"
	"tech-challenge-payment/internal/integration/order"
	"tech-challenge-payment/internal/repository"
	"tech-challenge-payment/internal/service"

	"github.com/sirupsen/logrus"
)

var (
	cfg = &config.Cfg
)

func main() {
	config.ParseFromFlags()
	restChannel := rest.NewPaymentChannel(service.NewPaymentService(repository.NewPaymentRepo(repository.NewMongo()), order.NewOrderService(http.DefaultClient)))
	if err := rest.New(restChannel).Start(); err != nil {
		logrus.Panic()
	}

}
