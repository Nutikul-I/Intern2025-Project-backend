package model

type (
	Permission struct {
		ID        int    `json:"id"`
		RoleID    int    `json:"role_id"`
		Module    string `json:"module"`
		CanView   bool   `json:"can_view"`
		CanCreate bool   `json:"can_create"`
		CanEdit   bool   `json:"can_edit"`
		CanDelete bool   `json:"can_delete"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		IsDeleted bool   `json:"is_deleted"`
		DeletedAt string `json:"deleted_at"`
	}

	PermissionPayload struct {
		ID        int    `json:"id"`
		RoleID    int    `json:"role_id"`
		Module    string `json:"module"`
		CanView   bool   `json:"can_view"`
		CanCreate bool   `json:"can_create"`
		CanEdit   bool   `json:"can_edit"`
		CanDelete bool   `json:"can_delete"`
	}

	PermissionResponse struct {
		TotalPages     int          `json:"total_pages"`
		PermissionList []Permission `json:"permission_list"`
	}

	CreatePermissionPayload struct {
		RoleID    int    `json:"role_id"`
		Module    string `json:"module"`
		CanView   bool   `json:"can_view"`
		CanCreate bool   `json:"can_create"`
		CanEdit   bool   `json:"can_edit"`
		CanDelete bool   `json:"can_delete"`
	}

	UpdatePermissionPayload struct {
		ID        int    `json:"id"`
		RoleID    int    `json:"role_id"`
		Module    string `json:"module"`
		CanView   bool   `json:"can_view"`
		CanCreate bool   `json:"can_create"`
		CanEdit   bool   `json:"can_edit"`
		CanDelete bool   `json:"can_delete"`
	}
)

// SQL Queries
var SQL_GET_PERMISSIONS = `
SELECT 
	id,
	role_id,
	module,
	can_view,
	can_create,
	can_edit,
	can_delete,
	created_at,
	updated_at,
	is_deleted,
	deleted_at
FROM permissions
WHERE (? = 0 OR id = ?)
AND is_deleted = 0
ORDER BY id DESC
LIMIT ?, ?`

var SQL_GET_TOTAL_PERMISSIONS = `
SELECT COUNT(*) AS TotalCount
FROM permissions
WHERE (? = 0 OR id = ?)
AND is_deleted = 0`

var SQL_CHECK_PERMISSION = `
SELECT 
	id,
	role_id,
	module
FROM permissions
WHERE id = ?
AND is_deleted = 0`

var SQL_CHECK_PERMISSION_DUPLICATE = `
SELECT 
	id,
	role_id,
	module
FROM permissions
WHERE role_id = ? AND module = ?
AND is_deleted = 0`

var SQL_CREATE_PERMISSION = `
INSERT INTO permissions (
	role_id,
	module,
	can_view,
	can_create,
	can_edit,
	can_delete,
	created_at,
	is_deleted
) VALUES (
	?,
	?,
	?,
	?,
	?,
	?,
	NOW(),
	0
)`

var SQL_UPDATE_PERMISSION = `
UPDATE permissions
SET
	role_id = ?,
	module = ?,
	can_view = ?,
	can_create = ?,
	can_edit = ?,
	can_delete = ?,
	updated_at = NOW()
WHERE id = ?
AND is_deleted = 0`

var SQL_DELETE_PERMISSION = `
UPDATE permissions
SET 
	is_deleted = 1,
	deleted_at = NOW()
WHERE id = ?`

var SQL_GET_PERMISSION_BY_ID = `
SELECT 
	id,
	role_id,
	module,
	can_view,
	can_create,
	can_edit,
	can_delete,
	created_at,
	updated_at,
	is_deleted,
	deleted_at
FROM permissions
WHERE id = ?
AND is_deleted = 0`
