package controller

import (
	"strconv"

	"payso-internal-api/model"
	"payso-internal-api/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type OrderController interface {
	List(*fiber.Ctx) error
	Detail(*fiber.Ctx) error
	UpdateStatus(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
}

type orderController struct{ svc service.OrderService }

func NewOrderController(s service.OrderService) OrderController { return &orderController{s} }

/* ---------- List ---------- */
func (ctl *orderController) List(c *fiber.Ctx) error {
	p, _ := strconv.Atoi(c.Query("page"))
	r, _ := strconv.Atoi(c.Query("row"))
	data, err := ctl.svc.List(c.Context(), p, r)
	if err != nil {
		log.Errorf("ListOrders: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot fetch orders"})
	}
	return c.JSON(data)
}

/* ---------- Detail ---------- */
func (ctl *orderController) Detail(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).SendString("invalid id")
	}
	d, err := ctl.svc.Detail(c.Context(), id)
	if err != nil {
		return c.Status(404).SendString("order not found")
	}
	return c.JSON(d)
}

/* ---------- Update Status ---------- */
func (ctl *orderController) UpdateStatus(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		return c.Status(400).SendString("invalid id")
	}

	var req model.OrderUpdate
	if err := c.BodyParser(&req); err != nil || req.StatusID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid payload"})
	}
	if err := ctl.svc.UpdateStatus(c.Context(), id, req.StatusID); err != nil {
		log.Errorf("UpdateStatus: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot update"})
	}
	return c.SendStatus(204)
}

/* ---------- Delete ---------- */
func (ctl *orderController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		return c.Status(400).SendString("invalid id")
	}

	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		log.Errorf("DeleteOrder: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot delete"})
	}
	return c.SendStatus(204)
}
