package repo

import (
	pb "github.com/asadbekGo/book-shop-order/genproto/order_service"
)

// UserStorageI ...
type OrderStorageI interface {
	CreateOrder(pb.Order) (pb.Order, error)
	GetOrder(id string) (pb.Order, error)
	ListOrders(page, limit int64) ([]*pb.Order, int64, error)
	UpdateOrder(pb.Order) (pb.Order, error)
	DeleteOrder(id string) error
}
