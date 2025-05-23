package controller

import (
	"fmt"
	"payso-internal-api/model"
	"payso-internal-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type EmployeeController interface {
	GetEmployees(c *fiber.Ctx) error
	CreateEmployee(c *fiber.Ctx) error
	UpdateEmployee(c *fiber.Ctx) error
	DeleteEmployee(c *fiber.Ctx) error
	GetEmployeeByID(c *fiber.Ctx) error
}

type employeeController struct {
	employeeService service.EmployeeService
}

func NewEmployeeController(employeeService service.EmployeeService) EmployeeController {
	return &employeeController{employeeService}
}

func (ctl *employeeController) GetEmployees(c *fiber.Ctx) error {
	log.Infof("==-- GetEmployees --==")

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

	res, err := ctl.employeeService.GetEmployeeService(requestID, Page, Row)
	if err != nil {
		log.Errorf("GetEmployees Error from service GetEmployeeService: %v", err)
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
		"message":      "GetEmployees",
		"data":         res.EmployeeList,
		"totalPages":   totalPages,
		"nextPage":     Page + 1,
		"previousPage": Page - 1,
		"currentPage":  Page,
	})
}

func (ctl *employeeController) CreateEmployee(c *fiber.Ctx) error {
	log.Infof("==-- CreateEmployee --==")

	rawBody := c.Body()
	log.Infof("Raw body received: %s", string(rawBody))

	// Debug: ดู Content-Type header
	contentType := c.Get("Content-Type")
	log.Infof("Content-Type: %s", contentType)

	var payload model.CreateEmployeePayload

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Create Employee Error parsing: %v", err)
		log.Errorf("Error type: %T", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body.",
			"data": map[string]interface{}{
				"error": err.Error(),
				"body":  string(rawBody),
			},
		})
	}

	// Debug: ดู parsed payload
	log.Infof("Parsed payload: %+v", payload)

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Create Employee Error parsing: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body.",
			"data":    err,
		})
	}

	res, err := ctl.employeeService.CreateEmployeeService(payload)
	if err != nil {
		log.Errorf("CreateEmployee Error from service CreateEmployeeService: %v", err)
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
		"message": "Employee created successfully",
		"data":    nil,
	})
}

func (ctl *employeeController) UpdateEmployee(c *fiber.Ctx) error {
	log.Infof("==-- UpdateEmployee --==")

	var payload model.UpdateEmployeePayload

	if err := c.BodyParser(&payload); err != nil {
		log.Errorf("Update Employee Error parsing: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body.",
			"data":    err,
		})
	}

	res, err := ctl.employeeService.UpdateEmployeeService(payload)
	if err != nil {
		log.Errorf("UpdateEmployee Error from service UpdateEmployeeService: %v", err)
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
		"message": "Employee updated successfully",
		"data":    nil,
	})
}

func (ctl *employeeController) DeleteEmployee(c *fiber.Ctx) error {
	log.Infof("==-- DeleteEmployee --==")

	requestID, err := strconv.Atoi(c.Query("id", "0"))
	if err != nil || requestID == 0 {
		log.Errorf("Delete Employee Error: Invalid ID")
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid employee ID",
			"data":    nil,
		})
	}

	res, err := ctl.employeeService.DeleteEmployeeService(requestID)
	if err != nil {
		log.Errorf("DeleteEmployee Error from service DeleteEmployeeService: %v", err)
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
		"message": "Employee deleted successfully",
		"data":    nil,
	})
}

func (ctl *employeeController) GetEmployeeByID(c *fiber.Ctx) error {
	log.Infof("==-- GetEmployeeByID --==")

	idParam := c.Query("id", "0")
	log.Infof("Query parameter 'id' value: %s", idParam)

	requestID, err := strconv.Atoi(idParam)
	if err != nil || requestID == 0 {
		log.Errorf("Get Employee Error: Invalid ID from query parameter: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid employee ID in query",
			"data":    nil,
		})
	}

	employee, err := ctl.employeeService.GetEmployeeByIDService(requestID)
	if err != nil {
		log.Errorf("GetEmployeeByID Error from service GetEmployeeByIDService: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Employee with ID %d not found", requestID),
			"data":    nil,
		})
	}

	if employee.ID == 0 {
		log.Errorf("Employee with ID %d returned has ID = 0", requestID)
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Employee with ID %d not found", requestID),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "GetEmployeeByID",
		"data":    employee,
	})
}
