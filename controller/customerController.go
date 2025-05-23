package controller

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"payso-internal-api/model"
	"payso-internal-api/service"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

/* ---------- helper: map error → status / message ---------- */

func translateError(err error) (int, string) {
	if err == nil {
		return fiber.StatusOK, ""
	}

	/* ------- sql.ErrNoRows -> 404 ------- */
	if errors.Is(err, sql.ErrNoRows) {
		return fiber.StatusNotFound, "customer not found"
	}

	/* ------- MySQL duplicate key ------- */
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		msg := strings.ToLower(me.Message)

		switch {
		case strings.Contains(msg, "email"):
			return fiber.StatusConflict, "email already exists"
		case strings.Contains(msg, "username"):
			return fiber.StatusConflict, "username already exists"
		case strings.Contains(msg, "national"):
			return fiber.StatusConflict, "national id already exists"
		case strings.Contains(msg, "phone"):
			return fiber.StatusConflict, "phone already exists"
		default:
			return fiber.StatusConflict, "duplicate data"
		}
	}

	/* ------- validation / business layer ------- */
	if errors.Is(err, service.ErrInvalidPayload) {
		return fiber.StatusBadRequest, err.Error()
	}

	/* ------- fallback 500 ------- */
	return fiber.StatusInternalServerError, "internal server error"
}

/* ---------- controller impl ---------- */

type CustomerController interface {
	GetCustomers(c *fiber.Ctx) error
	GetCustomer(c *fiber.Ctx) error
	CreateCustomer(c *fiber.Ctx) error
	UpdateCustomer(c *fiber.Ctx) error
	DeleteCustomer(c *fiber.Ctx) error
}

type customerController struct{ svc service.CustomerService }

func NewCustomerController(s service.CustomerService) CustomerController {
	return &customerController{s}
}

/* ---------- LIST ---------- */
func (ctl *customerController) GetCustomers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	row, _ := strconv.Atoi(c.Query("row"))

	list, err := ctl.svc.GetCustomers(c.Context(), page, row)
	if err != nil {
		code, msg := translateError(err)
		log.Errorf("GetCustomers: %v", err)
		return c.Status(code).JSON(fiber.Map{"message": msg})
	}
	return c.JSON(list)
}

/* ---------- DETAIL ---------- */
func (ctl *customerController) GetCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}
	detail, err := ctl.svc.GetCustomer(c.Context(), id)
	if err != nil {
		code, msg := translateError(err)
		return c.Status(code).JSON(fiber.Map{"message": msg})
	}
	return c.JSON(detail)
}

/* ---------- CREATE ---------- */
func (ctl *customerController) CreateCustomer(c *fiber.Ctx) error {
	var req model.CustomerCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
	}
	if pw := strings.TrimSpace(req.Password); pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password required"})
	}
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password too short (min 6)"})
	}

	id, err := ctl.svc.CreateCustomer(c.Context(), req)
	if err != nil {
		code, msg := translateError(err)
		log.Errorf("CreateCustomer: %v", err)
		return c.Status(code).JSON(fiber.Map{"message": msg})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"CustomerID": id})
}

/* ---------- UPDATE ---------- */
func (ctl *customerController) UpdateCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	var req model.CustomerCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload"})
	}
	if pw := strings.TrimSpace(req.Password); pw != "" && len(pw) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password too short (min 6)"})
	}

	if err := ctl.svc.UpdateCustomer(c.Context(), id, req); err != nil {
		code, msg := translateError(err)
		log.Errorf("UpdateCustomer: %v", err)
		return c.Status(code).JSON(fiber.Map{"message": msg})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

/* ---------- DELETE ---------- */
func (ctl *customerController) DeleteCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}

	if err := ctl.svc.DeleteCustomer(c.Context(), id); err != nil {
		code, msg := translateError(err)
		log.Errorf("DeleteCustomer: %v", err)
		return c.Status(code).JSON(fiber.Map{"message": msg})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
