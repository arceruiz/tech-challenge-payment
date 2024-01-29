package mocks

import "github.com/stretchr/testify/mock"

type OrderServiceMock struct {
	mock.Mock
}

func (m *OrderServiceMock) UpdateStatus(id, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}
