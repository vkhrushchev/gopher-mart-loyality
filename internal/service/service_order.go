package service

import (
	"context"
	"fmt"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

var (
	ErrOrderUploadedByAnotherUser = fmt.Errorf("order already uploaded by another user")
	ErrOrderWrongNumber           = fmt.Errorf("order wrong number")
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) PutOrder(ctx context.Context, orderNumber string) (bool, error) {
	log.Infow(
		"service: put order",
		"number", orderNumber)

	return false, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]dto.OrderDomain, error) {
	log.Infow(
		"service: get order list for user")

	return make([]dto.OrderDomain, 0), nil
}
