package controller

import (
	"payso-internal-api/model"
	"payso-internal-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type MerchantController interface {
	GetMerchant(c *fiber.Ctx) error
	CreateMerchant(c *fiber.Ctx) error
	DeleteMerchant(c *fiber.Ctx) error
}

type merchantController struct {
	merchantService service.MerchantService
}

func NewMerchantController(merchantService service.MerchantService) MerchantController {
	return &merchantController{merchantService}
}

func (ctl *merchantController) GetMerchant(c *fiber.Ctx) error {
	log.Infof("==-- GetMerchant --==")

	var Page int = 1
	var Row int = 50

	RequestMID := c.Query("MID")
	RequestPage, err := strconv.Atoi(c.Query("Page"))
	RequestRow, err := strconv.Atoi(c.Query("Row"))

	if RequestPage > 0 {
		Page = RequestPage
	}

	if RequestRow > 0 {
		Row = RequestRow
	}

	res, err := ctl.merchantService.GetMerchantService(RequestMID, Page, Row)
	if err != nil {
		log.Error("GetMerchant Error from service GetMerchant: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	totalPages := res.TotalPages / Row
	if res.TotalPages%Row != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       200,
		"message":      "GetMerchant",
		"data":         res.MerchantList,
		"totalPages":   totalPages,
		"nextPage":     Page + 1,
		"previousPage": Page - 1,
		"currentPage":  Page,
	})
}

func (ctl *merchantController) CreateMerchant(c *fiber.Ctx) error {
	log.Infof("==-- CreateMerchant --==")

	var payload model.CreateMerchantPayload

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Create Connection Type Error parsing")
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	res, err := ctl.merchantService.CreateMerchantService(payload, c.IP())
	if err != nil {
		log.Error("CreateMerchant Error from service CreateMerchant: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	if res.StatusCode == 400 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  400,
			"message": res.Message,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "CreateMerchant",
		"data":    nil,
	})
}

func (ctl *merchantController) DeleteMerchant(c *fiber.Ctx) error {
	log.Infof("==-- DeleteMerchant --==")

	ReqMasterMerchantID := c.Query("MasterMerchantID")
	ReqMerchantID := c.Query("MerchantID")

	res, err := ctl.merchantService.DeleteMerchantService(ReqMasterMerchantID, ReqMerchantID)
	if err != nil {
		log.Error("DeleteMerchant Error from service DeleteMerchant: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	if res.StatusCode == 400 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  400,
			"message": res.Message,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "DeleteMerchant",
		"data":    nil,
	})
}
