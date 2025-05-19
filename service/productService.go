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

/* --------- (optional) total count ถ้าต้องใช้หน้า UI ---------
func (s *productService) TotalProduct(ctx context.Context) (int, error) {
	return productRepository.GetTotalProductRepository(0)
}
---------------------------------------------------------------- */
