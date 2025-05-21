package handler

import (
	"context"
	"payso-internal-api/model"
	customerRepo "payso-internal-api/repository"
)

type CustomerHandler interface {
	GetCustomers(ctx context.Context, page, row int) ([]model.CustomerDetail, error)
	GetCustomer(ctx context.Context, id int64) (model.CustomerDetail, error)
	CreateCustomer(ctx context.Context, c model.CustomerCreate) (int64, error)
	UpdateCustomer(ctx context.Context, id int64, c model.CustomerCreate) error
	DeleteCustomer(ctx context.Context, id int64) error
}

type customerHandler struct{}

func NewCustomerHandler() CustomerHandler { return &customerHandler{} }

func (h *customerHandler) GetCustomers(ctx context.Context, page, row int) ([]model.CustomerDetail, error) {
	return customerRepo.GetCustomerList(ctx, page, row)
}
func (h *customerHandler) GetCustomer(ctx context.Context, id int64) (model.CustomerDetail, error) {
	return customerRepo.GetCustomerDetail(ctx, id)
}
func (h *customerHandler) CreateCustomer(ctx context.Context, c model.CustomerCreate) (int64, error) {
	return customerRepo.CreateCustomer(ctx, c)
}
func (h *customerHandler) UpdateCustomer(ctx context.Context, id int64, c model.CustomerCreate) error {
	return customerRepo.UpdateCustomer(ctx, id, c)
}
func (h *customerHandler) DeleteCustomer(ctx context.Context, id int64) error {
	return customerRepo.DeleteCustomer(ctx, id)
}
