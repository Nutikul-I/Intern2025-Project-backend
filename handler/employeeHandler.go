package handler

import (
	"fmt"
	"payso-internal-api/model"
	"payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

type EmployeeHandler interface {
	GetEmployees(id int, page int, row int) ([]model.Employee, error)
	GetTotalEmployees(id int) (int, error)
	CreateEmployee(payload model.CreateEmployeePayload) (model.UpdateResponse, error)
	UpdateEmployee(payload model.UpdateEmployeePayload) (model.UpdateResponse, error)
	DeleteEmployee(id int) (model.UpdateResponse, error)
	GetEmployeeByID(id int) (model.Employee, error)
}

type employeeHandler struct {
	// อาจมีตัวแปรสำหรับการกำหนดค่าเพิ่มเติมหรือ dependencies
}

func NewEmployeeHandler() EmployeeHandler {
	return &employeeHandler{}
}

func (h *employeeHandler) GetEmployees(id int, page int, row int) ([]model.Employee, error) {
	log.Infof("==-- Handler.GetEmployees --==")

	// ในที่นี้เราจะเรียกใช้งาน repository โดยตรง
	// แต่ในกรณีที่มีตรรกะเพิ่มเติม เช่น การคำนวณหรือปรับแต่งข้อมูล สามารถทำได้ที่นี่
	employees, err := repository.GetEmployeesRepository(id, page, row)
	if err != nil {
		log.Errorf("Error from repository.GetEmployeesRepository: %v", err)
		return nil, err
	}

	// อาจมีการปรับแต่งข้อมูลเพิ่มเติมหรือประมวลผลข้อมูลก่อนส่งกลับ
	// เช่น การเติมข้อมูลเพิ่มเติม หรือการกรองข้อมูลบางส่วน

	return employees, nil
}

func (h *employeeHandler) GetTotalEmployees(id int) (int, error) {
	log.Infof("==-- Handler.GetTotalEmployees --==")

	totalCount, err := repository.GetTotalEmployeesRepository(id)
	if err != nil {
		log.Errorf("Error from repository.GetTotalEmployeesRepository: %v", err)
		return 0, err
	}

	return totalCount, nil
}

func (h *employeeHandler) CreateEmployee(payload model.CreateEmployeePayload) (model.UpdateResponse, error) {
	log.Infof("==-- Handler.CreateEmployee --==")

	// อาจมีการตรวจสอบข้อมูลเพิ่มเติมก่อนส่งไปยัง repository
	// เช่น การตรวจสอบความถูกต้องของข้อมูล หรือการเตรียมข้อมูลเพิ่มเติม

	// ตรวจสอบความถูกต้องของข้อมูล
	if payload.UserID <= 0 {
		log.Errorf("Invalid user_id: %d", payload.UserID)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid user_id. user_id must be greater than 0",
		}, nil
	}

	if payload.RoleID <= 0 {
		log.Errorf("Invalid role_id: %d", payload.RoleID)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid role_id. role_id must be greater than 0",
		}, nil
	}

	response, err := repository.CreateEmployeeRepository(payload)
	if err != nil {
		log.Errorf("Error from repository.CreateEmployeeRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return response, nil
}

func (h *employeeHandler) UpdateEmployee(payload model.UpdateEmployeePayload) (model.UpdateResponse, error) {
	log.Infof("==-- Handler.UpdateEmployee --==")

	// ตรวจสอบความถูกต้องของข้อมูล
	if payload.ID <= 0 {
		log.Errorf("Invalid id: %d", payload.ID)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid id. ID must be greater than 0",
		}, nil
	}

	if payload.UserID <= 0 {
		log.Errorf("Invalid user_id: %d", payload.UserID)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid user_id. user_id must be greater than 0",
		}, nil
	}

	if payload.RoleID <= 0 {
		log.Errorf("Invalid role_id: %d", payload.RoleID)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid role_id. role_id must be greater than 0",
		}, nil
	}

	response, err := repository.UpdateEmployeeRepository(payload)
	if err != nil {
		log.Errorf("Error from repository.UpdateEmployeeRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return response, nil
}

func (h *employeeHandler) DeleteEmployee(id int) (model.UpdateResponse, error) {
	log.Infof("==-- Handler.DeleteEmployee --==")

	// ตรวจสอบความถูกต้องของข้อมูล
	if id <= 0 {
		log.Errorf("Invalid id: %d", id)
		return model.UpdateResponse{
			StatusCode: 400,
			Message:    "Invalid id. ID must be greater than 0",
		}, nil
	}

	response, err := repository.DeleteEmployeeRepository(id)
	if err != nil {
		log.Errorf("Error from repository.DeleteEmployeeRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return response, nil
}

func (h *employeeHandler) GetEmployeeByID(id int) (model.Employee, error) {
	log.Infof("==-- Handler.GetEmployeeByID --==")

	// ตรวจสอบความถูกต้องของข้อมูล
	if id <= 0 {
		log.Errorf("Invalid id: %d", id)
		return model.Employee{}, fmt.Errorf("invalid id: %d", id)
	}

	employee, err := repository.GetEmployeeByIDRepository(id)
	if err != nil {
		log.Errorf("Error from repository.GetEmployeeByIDRepository: %v", err)
		return model.Employee{}, err
	}

	// อาจมีการปรับแต่งข้อมูลเพิ่มเติมหรือการเติมข้อมูลจากแหล่งอื่น
	// เช่น เพิ่มข้อมูลจาก user หรือ role ถ้าต้องการ

	return employee, nil
}
