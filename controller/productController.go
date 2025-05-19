package controller

import (
	"strconv"

	"payso-internal-api/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type ProductController interface {
	GetProducts(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
}

type productController struct {
	svc service.ProductService
}

func NewProductController(svc service.ProductService) ProductController {
	return &productController{svc}
}

func (ctl *productController) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	row, _ := strconv.Atoi(c.Query("row"))

	products, err := ctl.svc.GetProducts(c.Context(), page, row)
	if err != nil {
		log.Errorf("GetProducts: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "cannot fetch products"})
	}
	return c.JSON(products)
}

func (ctl *productController) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}
	prod, err := ctl.svc.GetProductDetail(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("product not found")
	}
	return c.JSON(prod)
}
