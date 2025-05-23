package model

type (
	Employee struct {
		ID     int `json:"id" db:"id"`
		UserID int `json:"user_id" db:"user_id"`
		RoleID int `json:"role_id" db:"role_id"`
	}

	EmployeePayload struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
		RoleID int `json:"role_id"`
	}

	EmployeeResponse struct {
		TotalPages   int        `json:"total_pages"`
		EmployeeList []Employee `json:"employee_list"`
	}

	CreateEmployeePayload struct {
		UserID int `json:"user_id"`
		RoleID int `json:"role_id"`
	}

	UpdateEmployeePayload struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
		RoleID int `json:"role_id"`
	}
)

// SQL Queries
var SQL_GET_EMPLOYEES = `
SELECT 
	id,
	user_id,
	role_id
FROM employees
WHERE (? = 0 OR id = ?)
ORDER BY id DESC
LIMIT ?, ?`

var SQL_GET_TOTAL_EMPLOYEES = `
SELECT COUNT(*) AS TotalCount
FROM employees
WHERE (? = 0 OR id = ?)`

var SQL_CHECK_EMPLOYEE = `
SELECT 
	id,
	user_id,
	role_id
FROM employees
WHERE id = ?`

var SQL_CHECK_EMPLOYEE_DUPLICATE = `
SELECT 
	id,
	user_id,
	role_id
FROM employees
WHERE user_id = ?`

var SQL_CREATE_EMPLOYEE = `
INSERT INTO employees (
	user_id,
	role_id
) VALUES (
	?,
	?
)`

var SQL_UPDATE_EMPLOYEE = `
UPDATE employees
SET
	user_id = ?,
	role_id = ?
WHERE id = ?`

var SQL_DELETE_EMPLOYEE = `
DELETE FROM employees
WHERE id = ?`

var SQL_GET_EMPLOYEE_BY_ID = `
SELECT 
	id,
	user_id,
	role_id
FROM employees
WHERE id = ?`
