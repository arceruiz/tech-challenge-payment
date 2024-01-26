package order

import (
	"errors"
	"net/http"
	"strconv"
	"tech-challenge-payment/internal/config"
)

const (
	orderEndpoint = "/api/order/"
)

var (
	cfg = &config.Cfg
)

type orderService struct {
	httpClient *http.Client
}

func NewOrderService(client *http.Client) OrderService {
	return &orderService{
		httpClient: client,
	}
}

func (s *orderService) UpdateStatus(id, status string) error {

	request, err := http.NewRequest("PATCH", cfg.Server.OrderServiceHost+orderEndpoint, nil)
	if err != nil {
		return err
	}
	q := request.URL.Query()
	q.Add("id", id)
	q.Add("status", status)
	request.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("order integration error, code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
