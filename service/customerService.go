package service

import (
	"context"
	"payso-internal-api/handler"
	"payso-internal-api/model"

	log "github.com/sirupsen/logrus"
)

type CustomerService interface {
	GetCustomers(ctx context.Context, page, row int) ([]model.CustomerDetail, error)
	GetCustomer(ctx context.Context, id int64) (model.CustomerDetail, error)
	CreateCustomer(ctx context.Context, c model.CustomerCreate) (int64, error)
	UpdateCustomer(ctx context.Context, id int64, c model.CustomerCreate) error
	DeleteCustomer(ctx context.Context, id int64) error
}

type customerService struct{ h handler.CustomerHandler }

func NewCustomerService(h handler.CustomerHandler) CustomerService { return &customerService{h} }

func (s *customerService) GetCustomers(ctx context.Context, p, r int) ([]model.CustomerDetail, error) {
	log.Info("==-- GetCustomers Service --==")
	return s.h.GetCustomers(ctx, p, r)
}
func (s *customerService) GetCustomer(ctx context.Context, id int64) (model.CustomerDetail, error) {
	log.Info("==-- GetCustomer Service --==")
	return s.h.GetCustomer(ctx, id)
}
func (s *customerService) CreateCustomer(ctx context.Context, c model.CustomerCreate) (int64, error) {
	log.Info("==-- CreateCustomer Service --==")
	return s.h.CreateCustomer(ctx, c)
}
func (s *customerService) UpdateCustomer(ctx context.Context, id int64, c model.CustomerCreate) error {
	log.Info("==-- UpdateCustomer Service --==")
	return s.h.UpdateCustomer(ctx, id, c)
}
func (s *customerService) DeleteCustomer(ctx context.Context, id int64) error {
	log.Info("==-- DeleteCustomer Service --==")
	return s.h.DeleteCustomer(ctx, id)
}
