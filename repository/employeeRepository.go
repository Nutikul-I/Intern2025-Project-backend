package repository

import (
	"context"
	"fmt"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

func GetEmployeesRepository(id int, page int, row int) ([]model.Employee, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return nil, err
	}

	var Offset int = 0
	if page > 0 {
		Offset = (page - 1) * row
	}

	rows, err := conn.QueryContext(ctx, model.SQL_GET_EMPLOYEES,
		id,     // แทน WHERE ((@ID = 0) OR (id = @ID))
		id,     // ต้องระบุซ้ำเพราะมีการใช้ ? สองครั้งในเงื่อนไข
		Offset, // แทน OFFSET @Offset
		row)    // แทน ROWS FETCH NEXT @Row
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var EmployeeList []model.Employee
	err = scan.Rows(&EmployeeList, rows)
	if err != nil {
		log.Errorf("Error scanning rows: %v", err)
		return EmployeeList, err
	}

	return EmployeeList, nil
}

func GetTotalEmployeesRepository(id int) (int, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return 0, err
	}

	rows, err := conn.QueryContext(ctx, model.SQL_GET_TOTAL_EMPLOYEES, id, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return 0, err
	}
	defer rows.Close()

	TotalCount := 0
	err = scan.Row(&TotalCount, rows)
	if err != nil {
		log.Errorf("Error scan row : %v", err)
		return 0, err
	}

	return TotalCount, nil
}

func CreateEmployeeRepository(body model.CreateEmployeePayload) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	// Check if employee already exists for this user_id
	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_EMPLOYEE_DUPLICATE,
		body.UserID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingEmployee model.Employee
	err = scan.Row(&existingEmployee, rows_check)
	if err == nil {
		// Found existing employee, return error
		return model.UpdateResponse{StatusCode: 400, Message: "Employee already exists for this user"}, nil
	}

	// Insert new employee
	_, err = conn.ExecContext(ctx, model.SQL_CREATE_EMPLOYEE,
		body.UserID,
		body.RoleID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Created employee successfully"}, nil
}

func UpdateEmployeeRepository(body model.UpdateEmployeePayload) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	// Check if employee exists
	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_EMPLOYEE, body.ID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingEmployee model.Employee
	err = scan.Row(&existingEmployee, rows_check)
	if err != nil {
		// Employee not found
		return model.UpdateResponse{StatusCode: 404, Message: "Employee not found"}, nil
	}

	// Check if updated user_id exists on another employee
	if existingEmployee.UserID != body.UserID {
		rows_dup, err := conn.QueryContext(ctx, model.SQL_CHECK_EMPLOYEE_DUPLICATE,
			body.UserID)
		if err != nil {
			log.Errorf("Error executing query: %v", err)
			return model.UpdateResponse{}, err
		}
		defer rows_dup.Close()

		var dupEmployee model.Employee
		err = scan.Row(&dupEmployee, rows_dup)
		if err == nil && dupEmployee.ID != body.ID {
			// Found duplicate employee
			return model.UpdateResponse{StatusCode: 400, Message: "Employee already exists for this user"}, nil
		}
	}

	// Update employee
	_, err = conn.ExecContext(ctx, model.SQL_UPDATE_EMPLOYEE,
		body.UserID,
		body.RoleID,
		body.ID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Updated employee successfully"}, nil
}

func DeleteEmployeeRepository(id int) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	// Check if employee exists
	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_EMPLOYEE, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingEmployee model.Employee
	err = scan.Row(&existingEmployee, rows_check)
	if err != nil {
		// Employee not found
		return model.UpdateResponse{StatusCode: 404, Message: "Employee not found"}, nil
	}

	// Perform soft delete
	_, err = conn.ExecContext(ctx, model.SQL_DELETE_EMPLOYEE, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Deleted employee successfully"}, nil
}

func GetEmployeeByIDRepository(id int) (model.Employee, error) {
	conn := ConnectDB()
	ctx := context.Background()

	log.Infof("GetEmployeeByIDRepository: Searching for ID = %d", id)

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.Employee{}, err
	}

	// เพิ่ม dump sql และ parameters เพื่อตรวจสอบ
	log.Infof("SQL Query: %s, Parameters: [%d]", model.SQL_GET_EMPLOYEE_BY_ID, id)

	rows, err := conn.QueryContext(ctx, model.SQL_GET_EMPLOYEE_BY_ID, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.Employee{}, err
	}
	defer rows.Close()

	// Check if rows are empty
	hasRows := rows.Next()
	if !hasRows {
		log.Errorf("No rows found for ID = %d", id)
		return model.Employee{}, fmt.Errorf("employee with ID = %d not found", id)
	}

	// Reset rows
	rows.Close()
	rows, err = conn.QueryContext(ctx, model.SQL_GET_EMPLOYEE_BY_ID, id)
	if err != nil {
		log.Errorf("Error re-executing query: %v", err)
		return model.Employee{}, err
	}
	defer rows.Close()

	var employee model.Employee

	// แทนที่จะใช้ scan.Row, ลองใช้ rows.Scan โดยตรง
	if rows.Next() {
		err = rows.Scan(
			&employee.ID,
			&employee.UserID,
			&employee.RoleID,
			// ลบ fields ที่ไม่มีในตารางออก
		)
		if err != nil {
			log.Errorf("Error scanning row: %v", err)
			return model.Employee{}, err
		}
	}

	// ตรวจสอบข้อมูลที่ได้รับ
	log.Infof("Retrieved employee data: %+v", employee)

	if employee.ID == 0 {
		log.Errorf("Employee ID is 0 despite selecting ID = %d", id)
		return model.Employee{}, fmt.Errorf("invalid employee data retrieved for ID = %d", id)
	}

	return employee, nil
}
