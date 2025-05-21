package repository

import (
	"context"
	"fmt"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

func GetPermissionsRepository(id int, page int, row int) ([]model.Permission, error) {
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

	rows, err := conn.QueryContext(ctx, model.SQL_GET_PERMISSIONS,
		id,     // แทน WHERE ((@ID = 0) OR (id = @ID))
		id,     // ต้องระบุซ้ำเพราะมีการใช้ ? สองครั้งในเงื่อนไข
		Offset, // แทน OFFSET @Offset
		row)    // แทน ROWS FETCH NEXT @Row
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var PermissionList []model.Permission
	err = scan.Rows(&PermissionList, rows)
	if err != nil {
		log.Errorf("Error scanning rows: %v", err)
		return PermissionList, err
	}

	return PermissionList, nil
}

func GetTotalPermissionsRepository(id int) (int, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return 0, err
	}

	rows, err := conn.QueryContext(ctx, model.SQL_GET_TOTAL_PERMISSIONS, id, id)
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

func CreatePermissionRepository(body model.CreatePermissionPayload) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_PERMISSION_DUPLICATE,
		body.RoleID,
		body.Module)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingPermission model.Permission
	err = scan.Row(&existingPermission, rows_check)
	if err == nil {

		return model.UpdateResponse{StatusCode: 400, Message: "Permission already exists for this role and module"}, nil
	}

	_, err = conn.ExecContext(ctx, model.SQL_CREATE_PERMISSION,
		body.RoleID,
		body.Module,
		body.CanView,
		body.CanCreate,
		body.CanEdit,
		body.CanDelete)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Created permission successfully"}, nil
}

func UpdatePermissionRepository(body model.UpdatePermissionPayload) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_PERMISSION, body.ID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingPermission model.Permission
	err = scan.Row(&existingPermission, rows_check)
	if err != nil {

		return model.UpdateResponse{StatusCode: 404, Message: "Permission not found"}, nil
	}

	if existingPermission.RoleID != body.RoleID || existingPermission.Module != body.Module {
		rows_dup, err := conn.QueryContext(ctx, model.SQL_CHECK_PERMISSION_DUPLICATE,
			body.RoleID,
			body.Module)
		if err != nil {
			log.Errorf("Error executing query: %v", err)
			return model.UpdateResponse{}, err
		}
		defer rows_dup.Close()

		var dupPermission model.Permission
		err = scan.Row(&dupPermission, rows_dup)
		if err == nil && dupPermission.ID != body.ID {
			// Found duplicate permission
			return model.UpdateResponse{StatusCode: 400, Message: "Permission already exists for this role and module"}, nil
		}
	}

	_, err = conn.ExecContext(ctx, model.SQL_UPDATE_PERMISSION,
		body.RoleID,
		body.Module,
		body.CanView,
		body.CanCreate,
		body.CanEdit,
		body.CanDelete,
		body.ID)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Updated permission successfully"}, nil
}

func DeletePermissionRepository(id int) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	rows_check, err := conn.QueryContext(ctx, model.SQL_CHECK_PERMISSION, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var existingPermission model.Permission
	err = scan.Row(&existingPermission, rows_check)
	if err != nil {

		return model.UpdateResponse{StatusCode: 404, Message: "Permission not found"}, nil
	}

	_, err = conn.ExecContext(ctx, model.SQL_DELETE_PERMISSION, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "Deleted permission successfully"}, nil
}

func GetPermissionByIDRepository(id int) (model.Permission, error) {
	conn := ConnectDB()
	ctx := context.Background()

	log.Infof("GetPermissionByIDRepository: Searching for ID = %d", id)

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.Permission{}, err
	}

	log.Infof("SQL Query: %s, Parameters: [%d]", model.SQL_GET_PERMISSION_BY_ID, id)

	rows, err := conn.QueryContext(ctx, model.SQL_GET_PERMISSION_BY_ID, id)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.Permission{}, err
	}
	defer rows.Close()

	hasRows := rows.Next()
	if !hasRows {
		log.Errorf("No rows found for ID = %d", id)
		return model.Permission{}, fmt.Errorf("permission with ID = %d not found", id)
	}

	rows.Close()
	rows, err = conn.QueryContext(ctx, model.SQL_GET_PERMISSION_BY_ID, id)
	if err != nil {
		log.Errorf("Error re-executing query: %v", err)
		return model.Permission{}, err
	}
	defer rows.Close()

	var permission model.Permission

	if rows.Next() {
		err = rows.Scan(
			&permission.ID,
			&permission.RoleID,
			&permission.Module,
			&permission.CanView,
			&permission.CanCreate,
			&permission.CanEdit,
			&permission.CanDelete,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.IsDeleted,
			&permission.DeletedAt,
		)
		if err != nil {
			log.Errorf("Error scanning row: %v", err)
			return model.Permission{}, err
		}
	} else {
		log.Errorf("No data found after resetting rows cursor")
		return model.Permission{}, fmt.Errorf("no data found for permission ID = %d", id)
	}

	log.Infof("Retrieved permission data: %+v", permission)

	if permission.ID == 0 {
		log.Errorf("Permission ID is 0 despite selecting ID = %d", id)
		return model.Permission{}, fmt.Errorf("invalid permission data retrieved for ID = %d", id)
	}

	return permission, nil
}
