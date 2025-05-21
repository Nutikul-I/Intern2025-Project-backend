package handler

import (
	"context"
	"payso-internal-api/model"
	productRepository "payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

// ProductHandler ทำหน้าที่คั่นกลางระหว่าง service ↔ repository
type ProductHandler interface {
	// page, row = 1-based / จำนวนต่อหน้า
	GetProducts(ctx context.Context, page, row int) ([]model.ProductPayload, error)
	GetProduct(ctx context.Context, id int64) (model.ProductDetail, error)
	CreateProduct(ctx context.Context, p model.ProductCreate) (int64, error)
	UpdateProduct(ctx context.Context, pid int64, p model.ProductCreate) error
	DeleteProduct(ctx context.Context, pid int64) error
}

type productHandler struct{}

func NewProductHandler() ProductHandler {
	return &productHandler{}
}

/* ---------- list ---------- */
func (h *productHandler) GetProducts(ctx context.Context, page, row int) ([]model.ProductPayload, error) {
	list, err := productRepository.GetProductRepository(0 /* pid 0 = ดึงทั้งหมด */, page, row)
	if err != nil {
		log.Errorf("ProductRepository list error: %v", err)
	}
	return list, err
}

/* ---------- detail ---------- */
func (h *productHandler) GetProduct(ctx context.Context, id int64) (model.ProductDetail, error) {
	detail, err := productRepository.GetProductDetailRepository(id)
	if err != nil {
		log.Errorf("ProductRepository detail error: %v", err)
	}
	return detail, err
}

func (h *productHandler) CreateProduct(ctx context.Context, p model.ProductCreate) (int64, error) {
	return productRepository.CreateProductRepository(ctx, p)
}

func (h *productHandler) UpdateProduct(ctx context.Context, pid int64, p model.ProductCreate) error {
	return productRepository.UpdateProductRepository(ctx, pid, p)
}
func (h *productHandler) DeleteProduct(ctx context.Context, pid int64) error {
	return productRepository.DeleteProductRepository(ctx, pid)
}
