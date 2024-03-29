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

func New() rest {
	return rest{
		payment: NewPaymentChannel(),
	}
}

func (r rest) Start() error {
	router := echo.New()

	router.Use(middlewares.Logger)

	mainGroup := router.Group("/api")
	mainGroup.GET("/healthz", r.payment.HealthCheck)
	paymentGroup := mainGroup.Group("/payment")
	paymentGroup.Use(middlewares.Authorization)
	r.payment.RegisterGroup(paymentGroup)

	return router.Start(":" + cfg.Server.Port)
}
