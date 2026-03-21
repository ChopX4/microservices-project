package order

import "github.com/ChopX4/raketka/order/internal/repository"

type service struct {
	orderRepository repository.OrderRepository
}

func NewService(orderRepository repository.OrderRepository) *service {
	return &service{
		orderRepository: orderRepository,
	}
}
