package service

import (
	"context"
	"payso-internal-api/handler"
	"payso-internal-api/model"

	log "github.com/sirupsen/logrus"
)

type OrderService interface {
	List(ctx context.Context, p, r int) (model.OrderPagination, error)
	Detail(ctx context.Context, id int64) (model.OrderDetail, error)
	UpdateStatus(ctx context.Context, id, statusID int64) error
	Delete(ctx context.Context, id int64) error
}

type orderService struct{ h handler.OrderHandler }

func NewOrderService(h handler.OrderHandler) OrderService { return &orderService{h} }

func (s *orderService) List(ctx context.Context, p, r int) (model.OrderPagination, error) {
	log.Info("==-- Order List --==")
	return s.h.List(ctx, p, r)
}
func (s *orderService) Detail(ctx context.Context, id int64) (model.OrderDetail, error) {
	log.Info("==-- Order Detail --==")
	return s.h.Detail(ctx, id)
}
func (s *orderService) UpdateStatus(ctx context.Context, id, statusID int64) error {
	log.Info("==-- Update Order Status --==")
	return s.h.UpdateStatus(ctx, id, statusID)
}
func (s *orderService) Delete(ctx context.Context, id int64) error {
	log.Info("==-- Delete Order --==")
	return s.h.Delete(ctx, id)
}
