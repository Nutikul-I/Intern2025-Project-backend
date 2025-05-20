// controller/permissionController.go
package controller

import (
	"payso-internal-api/model"
	"payso-internal-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type PermissionController interface {
	GetPermissions(c *fiber.Ctx) error
	CreatePermission(c *fiber.Ctx) error
	UpdatePermission(c *fiber.Ctx) error
	DeletePermission(c *fiber.Ctx) error
	GetPermissionByID(c *fiber.Ctx) error
}

type permissionController struct {
	permissionService service.PermissionService
}

func NewPermissionController(permissionService service.PermissionService) PermissionController {
	return &permissionController{permissionService}
}

func (ctl *permissionController) GetPermissions(c *fiber.Ctx) error {
	log.Infof("==-- GetPermissions --==")

	var Page int = 1
	var Row int = 50

	requestID, _ := strconv.Atoi(c.Query("id", "0"))
	requestPage, err := strconv.Atoi(c.Query("page", "1"))
	requestRow, err := strconv.Atoi(c.Query("row", "50"))

	if requestPage > 0 {
		Page = requestPage
	}

	if requestRow > 0 {
		Row = requestRow
	}

	res, err := ctl.permissionService.GetPermissionService(requestID, Page, Row)
	if err != nil {
		log.Errorf("GetPermissions Error from service GetPermissionService: %v", err)
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
		"message":      "GetPermissions",
		"data":         res.PermissionList,
		"totalPages":   totalPages,
		"nextPage":     Page + 1,
		"previousPage": Page - 1,
		"currentPage":  Page,
	})
}

func (ctl *permissionController) CreatePermission(c *fiber.Ctx) error {
	log.Infof("==-- CreatePermission --==")

	var payload model.CreatePermissionPayload

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Create Permission Error parsing: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body.",
			"data":    err,
		})
	}

	res, err := ctl.permissionService.CreatePermissionService(payload)
	if err != nil {
		log.Errorf("CreatePermission Error from service CreatePermissionService: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	if res.StatusCode != 200 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  res.StatusCode,
			"message": res.Message,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "Permission created successfully",
		"data":    nil,
	})
}

func (ctl *permissionController) UpdatePermission(c *fiber.Ctx) error {
	log.Infof("==-- UpdatePermission --==")

	var payload model.UpdatePermissionPayload

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Update Permission Error parsing: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body.",
			"data":    err,
		})
	}

	res, err := ctl.permissionService.UpdatePermissionService(payload)
	if err != nil {
		log.Errorf("UpdatePermission Error from service UpdatePermissionService: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	if res.StatusCode != 200 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  res.StatusCode,
			"message": res.Message,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "Permission updated successfully",
		"data":    nil,
	})
}

func (ctl *permissionController) DeletePermission(c *fiber.Ctx) error {
	log.Infof("==-- DeletePermission --==")

	requestID, err := strconv.Atoi(c.Query("id", "0"))
	if err != nil || requestID == 0 {
		log.Errorf("Delete Permission Error: Invalid ID")
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid permission ID",
			"data":    nil,
		})
	}

	res, err := ctl.permissionService.DeletePermissionService(requestID)
	if err != nil {
		log.Errorf("DeletePermission Error from service DeletePermissionService: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	if res.StatusCode != 200 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  res.StatusCode,
			"message": res.Message,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "Permission deleted successfully",
		"data":    nil,
	})
}

func (ctl *permissionController) GetPermissionByID(c *fiber.Ctx) error {
	log.Infof("==-- GetPermissionByID --==")

	requestID, err := strconv.Atoi(c.Query("id", "0"))
	if err != nil || requestID == 0 {
		log.Errorf("Get Permission Error: Invalid ID")
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid permission ID",
			"data":    nil,
		})
	}

	permission, err := ctl.permissionService.GetPermissionByIDService(requestID)
	if err != nil {
		log.Errorf("GetPermissionByID Error from service GetPermissionByIDService: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "API Failed.",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "GetPermissionByID",
		"data":    permission,
	})
}
