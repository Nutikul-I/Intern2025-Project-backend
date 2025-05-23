package service

import (
	"fmt"
	"payso-internal-api/handler"
	"payso-internal-api/model"

	log "github.com/sirupsen/logrus"
)

type EmployeeService interface {
	GetEmployeeService(id int, page int, row int) (model.EmployeeResponse, error)
	CreateEmployeeService(body model.CreateEmployeePayload) (model.UpdateResponse, error)
	UpdateEmployeeService(body model.UpdateEmployeePayload) (model.UpdateResponse, error)
	DeleteEmployeeService(id int) (model.UpdateResponse, error)
	GetEmployeeByIDService(id int) (model.Employee, error)
}

type employeeService struct {
	employeeHandler handler.EmployeeHandler // เพิ่ม field นี้
}

// แก้ไขฟังก์ชันนี้ให้รับ parameter
func NewEmployeeService(employeeHandler handler.EmployeeHandler) EmployeeService {
	return &employeeService{employeeHandler: employeeHandler}
}

func (s *employeeService) GetEmployeeService(id int, page int, row int) (model.EmployeeResponse, error) {
	log.Infof("==-- GetEmployeeService --==")

	var err error
	var EmployeeList []model.Employee

	// ใช้ s.employeeHandler แทนการเรียก repository โดยตรง
	EmployeeList, err = s.employeeHandler.GetEmployees(id, page, row)
	if err != nil {
		log.Errorf("Error from GetEmployees handler: %v", err)
		return model.EmployeeResponse{}, err
	}

	TotalPages, err := s.employeeHandler.GetTotalEmployees(id)
	if err != nil {
		log.Errorf("Error from GetTotalEmployees handler: %v", err)
		return model.EmployeeResponse{}, err
	}

	EmployeeResponse := model.EmployeeResponse{
		TotalPages:   TotalPages,
		EmployeeList: EmployeeList,
	}

	return EmployeeResponse, err
}

func (s *employeeService) CreateEmployeeService(body model.CreateEmployeePayload) (model.UpdateResponse, error) {
	log.Infof("==-- CreateEmployeeService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = s.employeeHandler.CreateEmployee(body)
	if err != nil {
		log.Errorf("Error from CreateEmployee handler: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *employeeService) UpdateEmployeeService(body model.UpdateEmployeePayload) (model.UpdateResponse, error) {
	log.Infof("==-- UpdateEmployeeService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = s.employeeHandler.UpdateEmployee(body)
	if err != nil {
		log.Errorf("Error from UpdateEmployee handler: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *employeeService) DeleteEmployeeService(id int) (model.UpdateResponse, error) {
	log.Infof("==-- DeleteEmployeeService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = s.employeeHandler.DeleteEmployee(id)
	if err != nil {
		log.Errorf("Error from DeleteEmployee handler: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *employeeService) GetEmployeeByIDService(id int) (model.Employee, error) {
	log.Infof("==-- GetEmployeeByIDService --== ID: %d", id)

	var err error
	var employee model.Employee

	employee, err = s.employeeHandler.GetEmployeeByID(id)
	if err != nil {
		log.Errorf("Error from GetEmployeeByID handler: %v", err)
		return model.Employee{}, err
	}

	// เพิ่มการตรวจสอบว่า ID ที่ได้กลับมาเป็น 0 หรือไม่
	if employee.ID == 0 {
		log.Errorf("Employee with ID %d not found or has ID = 0", id)
		return model.Employee{}, fmt.Errorf("employee with ID %d not found or has ID = 0", id)
	}

	return employee, nil
}
