package controller

import (
	"strconv"

	"payso-internal-api/model"
	"payso-internal-api/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type ProductController interface {
	GetProducts(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
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
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	prod, err := ctl.svc.GetProductDetail(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("product not found")
	}
	return c.JSON(prod)
}

func (ctl *productController) CreateProduct(c *fiber.Ctx) error {
	var req model.ProductCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
	}
	id, err := ctl.svc.CreateProduct(c.Context(), req)
	if err != nil {
		log.Errorf("CreateProduct: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "cannot create"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ProductID": id})
}

/* ---------- UPDATE ---------- */
func (ctl *productController) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Query("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || pid <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	var req model.ProductCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
	}
	if err := ctl.svc.UpdateProduct(c.Context(), pid, req); err != nil {
		log.Errorf("UpdateProduct: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "cannot update"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

/* ---------- DELETE ---------- */
func (ctl *productController) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Query("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || pid <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	if err := ctl.svc.DeleteProduct(c.Context(), pid); err != nil {
		log.Errorf("DeleteProduct: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "cannot delete"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
