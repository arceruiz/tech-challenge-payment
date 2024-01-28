package rest

import (
	"fmt"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/service"

	"net/http"

	"github.com/labstack/echo/v4"
)

type Payment interface {
	RegisterGroup(g *echo.Group)
	Callback(c echo.Context) error
	GetByID(c echo.Context) error
	Create(c echo.Context) error
}

type payment struct {
	paymentSvc service.PaymentService
}

func NewPaymentChannel(paymentService service.PaymentService) Payment {
	return &payment{
		paymentSvc: paymentService,
	}
}

func (p *payment) RegisterGroup(g *echo.Group) {
	g.GET("/:id", p.GetByID)
	g.GET("", p.GetAll)
	g.POST("/callback", p.Callback)
	g.POST("/", p.Create)
}

func (p *payment) Create(c echo.Context) error {
	var paymentRequest PaymentRequest

	if err := c.Bind(&paymentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "Invalid request body",
		})
	}
	payment, err := p.paymentSvc.Create(c.Request().Context(), paymentRequest.toCanonical())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error creating payment")
	}

	return c.JSON(http.StatusOK, payment)
}

func (p *payment) GetByID(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "missing id query param",
		})
	}

	payment, err := p.paymentSvc.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "error searching payment")
	}

	return c.JSON(http.StatusOK, payment)
}

func (p *payment) GetAll(c echo.Context) error {
	payments, err := p.paymentSvc.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusNotFound, "error searching payment")
	}

	return c.JSON(http.StatusOK, payments)
}

func (p *payment) Callback(c echo.Context) error {

	var callback PaymentCallback
	if err := c.Bind(&callback); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	if _, ok := canonical.MapPaymentStatus[callback.Status]; !ok {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid status").Error(),
		})
	}

	err := p.paymentSvc.Callback(c.Request().Context(), callback.PaymentID, canonical.MapPaymentStatus[callback.Status])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error processing callback")
	}

	return c.JSON(http.StatusOK, nil)
}
