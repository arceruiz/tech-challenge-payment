package rest

import (
	"tech-challenge-payment/internal/config"
	"tech-challenge-payment/internal/middlewares"

	"github.com/labstack/echo/v4"
)

var (
	cfg = &config.Cfg
)

type rest struct {
	payment Payment
}

func New(payment Payment) rest {
	return rest{
		payment: payment,
	}
}

func (r rest) Start() error {
	router := echo.New()

	router.Use(middlewares.Logger)
	router.Use(middlewares.Authorization)

	mainGroup := router.Group("/api")
	paymentGroup := mainGroup.Group("/payment")

	r.payment.RegisterGroup(paymentGroup)

	return router.Start(":" + cfg.Server.Port)
}
