package controller

import (
	"payso-internal-api/model"
	"payso-internal-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type CategoryController interface {
	GetCategory(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &categoryController{categoryService}
}

func (ctl *categoryController) GetCategory(c *fiber.Ctx) error {
	log.Infof("==-- GetCategory --==")

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

	res, err := ctl.categoryService.GetCategoryService(mid, page, row)
	if err != nil {
		log.Errorf("GetCategory Error from service: %v", err)
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
		"message":      "GetCategory",
		"data":         res.CategoryList,
		"totalPages":   totalPages,
		"nextPage":     page + 1,
		"previousPage": page - 1,
		"currentPage":  page,
	})
}

func (ctl *categoryController) CreateCategory(c *fiber.Ctx) error {
	log.Infof("==-- CreateCategory --==")

	var payload model.CreateCategory

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("CreateCategory Error parsing body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid payload.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.categoryService.CreateCategoryService(payload)
	if err != nil {
		log.Errorf("CreateCategory Error from service: %v", err)
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

func (ctl *categoryController) UpdateCategory(c *fiber.Ctx) error {
	log.Infof("==-- UpdateCategory --==")

	var payload model.UpdateCategory

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("UpdateCategory Error parsing body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid payload.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.categoryService.UpdateCategoryService(payload)
	if err != nil {
		log.Errorf("UpdateCategory Error from service: %v", err)
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

func (ctl *categoryController) DeleteCategory(c *fiber.Ctx) error {
	log.Infof("==-- DeleteCategory --==")

	categoryIDStr := c.Query("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Errorf("DeleteCategory invalid CategoryID: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid CategoryID.",
			"data":    err.Error(),
		})
	}

	res, err := ctl.categoryService.DeleteCategoryService(categoryID)
	if err != nil {
		log.Errorf("DeleteCategory Error from service: %v", err)
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
