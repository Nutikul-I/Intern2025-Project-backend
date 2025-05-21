package model

type (
	CategoryPlayload struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	CategoryPagination struct {
		TotalPages   int `json:"TotalPages"`
		CategoryList []CategoryPlayload
	}

	CreateCategory struct {
		Name string `json:"name"`
	}

	UpdateCategory struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
)

var SQL_GET_TOTAL_CATEGORY = `
SELECT COUNT(*) AS TotalCount
FROM categories 
WHERE is_deleted = FALSE
`

var SQL_GET_CATEGORY = `
SELECT 
    id, 
    name 
FROM categories 
WHERE is_deleted = FALSE 
ORDER BY id ASC
`

var SQL_CREATE_CATEGORY = `
INSERT INTO categories (name, created_at, updated_at, is_deleted)
VALUES (?, NOW(), NOW(), FALSE);
`

var SQL_UPDATE_CATEGORY = `
UPDATE categories
SET name = ?, updated_at = NOW()
WHERE id = ? AND is_deleted = FALSE;
`

var SQL_SOFT_DELETE_CATEGORY = `
UPDATE categories
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = ? AND is_deleted = FALSE;
`
