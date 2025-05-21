package controller

import (
	"strconv"
	"strings"

	"payso-internal-api/model"
	"payso-internal-api/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type CustomerController interface {
	GetCustomers(c *fiber.Ctx) error
	GetCustomer(c *fiber.Ctx) error
	CreateCustomer(c *fiber.Ctx) error
	UpdateCustomer(c *fiber.Ctx) error
	DeleteCustomer(c *fiber.Ctx) error
}

type customerController struct {
	svc service.CustomerService
}

func NewCustomerController(s service.CustomerService) CustomerController {
	return &customerController{s}
}

/* ---------- LIST ---------- */
func (ctl *customerController) GetCustomers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	row, _ := strconv.Atoi(c.Query("row"))

	list, err := ctl.svc.GetCustomers(c.Context(), page, row)
	if err != nil {
		log.Errorf("GetCustomers: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot fetch"})
	}
	return c.JSON(list)
}

/* ---------- DETAIL ---------- */
func (ctl *customerController) GetCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).SendString("invalid id")
	}
	detail, err := ctl.svc.GetCustomer(c.Context(), id)
	if err != nil {
		return c.Status(404).SendString("customer not found")
	}
	return c.JSON(detail)
}

/* ---------- CREATE ---------- */
func (ctl *customerController) CreateCustomer(c *fiber.Ctx) error {
	var req model.CustomerCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid payload"})
	}

	// ---- PASSWORD VALIDATION ----
	if pw := strings.TrimSpace(req.Password); pw == "" {
		return c.Status(400).JSON(fiber.Map{"message": "password required"})
	}
	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"message": "password too short (min 6)"})
	}

	id, err := ctl.svc.CreateCustomer(c.Context(), req)
	if err != nil {
		log.Errorf("CreateCustomer: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot create"})
	}
	return c.Status(201).JSON(fiber.Map{"CustomerID": id})
}

/* ---------- UPDATE ---------- */
func (ctl *customerController) UpdateCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).SendString("invalid id")
	}

	var req model.CustomerCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid payload"})
	}

	// ---- OPTIONAL PASSWORD ----
	if pw := strings.TrimSpace(req.Password); pw != "" && len(pw) < 6 {
		return c.Status(400).JSON(fiber.Map{"message": "password too short (min 6)"})
	}

	if err := ctl.svc.UpdateCustomer(c.Context(), id, req); err != nil {
		log.Errorf("UpdateCustomer: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot update"})
	}
	return c.SendStatus(204)
}

/* ---------- DELETE ---------- */
func (ctl *customerController) DeleteCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).SendString("invalid id")
	}

	if err := ctl.svc.DeleteCustomer(c.Context(), id); err != nil {
		log.Errorf("DeleteCustomer: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "cannot delete"})
	}
	return c.SendStatus(204)
}
