package controller

import (
	"payso-internal-api/model"
	"payso-internal-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type DiscountController interface {
	GetDiscount(c *fiber.Ctx) error
	CreateDiscount(c *fiber.Ctx) error
	UpdateDiscount(c *fiber.Ctx) error
	DeleteDiscount(c *fiber.Ctx) error
}

type discountController struct {
	discountService service.DiscountService
}

func NewDiscountController(discountService service.DiscountService) DiscountController {
	return &discountController{discountService}
}

func (ctl *discountController) GetDiscount(c *fiber.Ctx) error {
	log.Infof("==-- GetDiscount --==")

	page := 1
	row := 50

	if p := c.Query("Page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if r := c.Query("Row"); r != "" {
		if parsedRow, err := strconv.Atoi(r); err == nil && parsedRow > 0 {
			row = parsedRow
		}
	}

	mid := c.Query("MID")

	res, err := ctl.discountService.GetDiscountService(mid, page, row)
	if err != nil {
		log.Errorf("GetDiscount Error from service: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err.Error(),
		})
	}

	totalPages := res.TotalPages / row
	if res.TotalPages%row != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       200,
		"message":      "GetDiscount",
		"data":         res.DiscountList,
		"totalPages":   totalPages,
		"nextPage":     page + 1,
		"previousPage": page - 1,
		"currentPage":  page,
	})
}

func (ctl *discountController) CreateDiscount(c *fiber.Ctx) error {
	log.Infof("==-- CreateDiscount --==")

	var payload model.CreateDiscount

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("CreateDiscount Error parsing body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid payload.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.discountService.CreateDiscountService(payload)
	if err != nil {
		log.Errorf("CreateDiscount Error from service: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  res.StatusCode,
		"message": res.Message,
		"data":    nil,
	})
}

func (ctl *discountController) UpdateDiscount(c *fiber.Ctx) error {
	log.Infof("==-- UpdateDiscount --==")

	var payload model.UpdateDiscount

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("UpdateDiscount Error parsing body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid payload.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.discountService.UpdateDiscountService(payload)
	if err != nil {
		log.Errorf("UpdateDiscount Error from service: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  res.StatusCode,
		"message": res.Message,
		"data":    nil,
	})
}

func (ctl *discountController) DeleteDiscount(c *fiber.Ctx) error {
	log.Infof("==-- DeleteDiscount --==")

	discountIDStr := c.Query("id")
	discountID, err := strconv.Atoi(discountIDStr)
	if err != nil {
		log.Errorf("DeleteDiscount invalid DiscountID: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid DiscountID.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.discountService.DeleteDiscountService(discountID)
	if err != nil {
		log.Errorf("DeleteDiscount Error from service: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  res.StatusCode,
		"message": res.Message,
		"data":    nil,
	})
}
