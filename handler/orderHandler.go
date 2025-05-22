package handler

import (
	"context"

	"payso-internal-api/model"
	orderRepo "payso-internal-api/repository"
)

type OrderHandler interface {
	List(ctx context.Context, p, r int) (model.OrderPagination, error)
	Detail(ctx context.Context, id int64) (model.OrderDetail, error)
	UpdateStatus(ctx context.Context, id, statusID int64) error
	Delete(ctx context.Context, id int64) error
}

type orderHandler struct{}

func NewOrderHandler() OrderHandler { return &orderHandler{} }

func (h *orderHandler) List(ctx context.Context, p, r int) (model.OrderPagination, error) {
	return orderRepo.GetOrderList(ctx, p, r)
}
func (h *orderHandler) Detail(ctx context.Context, id int64) (model.OrderDetail, error) {
	return orderRepo.GetOrderDetail(ctx, id)
}
func (h *orderHandler) UpdateStatus(ctx context.Context, id, statusID int64) error {
	return orderRepo.UpdateOrderStatus(ctx, id, statusID)
}
func (h *orderHandler) Delete(ctx context.Context, id int64) error {
	return orderRepo.DeleteOrder(ctx, id)
}
