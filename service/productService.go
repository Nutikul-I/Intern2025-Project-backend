package service

import (
	"context"
	"payso-internal-api/handler"
	"payso-internal-api/model"

	log "github.com/sirupsen/logrus"
)

// -------------------------------- interface
type ProductService interface {
	GetProducts(ctx context.Context, page, row int) ([]model.ProductPayload, error)
	GetProductDetail(ctx context.Context, id int64) (model.ProductDetail, error)
	CreateProduct(ctx context.Context, p model.ProductCreate) (int64, error)
	UpdateProduct(ctx context.Context, pid int64, p model.ProductCreate) error
	DeleteProduct(ctx context.Context, pid int64) error
}
type productService struct {
	productHandler handler.ProductHandler
}

func NewProductService(h handler.ProductHandler) ProductService {
	return &productService{h}
}

/* -------- รายการ -------- */
func (s *productService) GetProducts(ctx context.Context, page, row int) ([]model.ProductPayload, error) {
	log.Infof("==-- GetProducts Service --==")
	return s.productHandler.GetProducts(ctx, page, row)
}

/* -------- รายละเอียด -------- */
func (s *productService) GetProductDetail(ctx context.Context, id int64) (model.ProductDetail, error) {
	log.Infof("==-- GetProductDetail Service --==")
	return s.productHandler.GetProduct(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, p model.ProductCreate) (int64, error) {
	log.Info("==-- CreateProduct Service --==")
	return s.productHandler.CreateProduct(ctx, p)
}
func (s *productService) UpdateProduct(ctx context.Context, pid int64, p model.ProductCreate) error {
	log.Info("==-- UpdateProduct Service --==")
	return s.productHandler.UpdateProduct(ctx, pid, p)
}
func (s *productService) DeleteProduct(ctx context.Context, pid int64) error {
	log.Info("==-- DeleteProduct Service --==")
	return s.productHandler.DeleteProduct(ctx, pid)
}
