package order

type OrderService interface {
	UpdateStatus(id, status string) error
}
