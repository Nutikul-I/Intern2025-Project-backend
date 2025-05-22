package model

type (
	DiscountPlayload struct {
		Id                int     `json:"id"`
		Code              string  `json:"code"`
		Amount            float64 `json:"amount"`
		TotalQuantity     int     `json:"totalQuantity" db:"total_quantity"`
		RemainingQuantity int     `json:"remainingQuantity" db:"remaining_quantity"`
	}

	DiscountPagination struct {
		TotalPages   int `json:"TotalPages"`
		DiscountList []DiscountPlayload
	}

	CreateDiscount struct {
		Code              string `json:"code"`
		Amount            int    `json:"amount"`
		TotalQuantity     int    `json:"totalQuantity" db:"total_quantity"`
		RemainingQuantity int    `json:"remainingQuantity" db:"remaining_quantity"`
	}

	UpdateDiscount struct {
		Id                int    `json:"id"`
		Code              string `json:"code"`
		Amount            int    `json:"amount"`
		TotalQuantity     int    `json:"totalQuantity" db:"total_quantity"`
		RemainingQuantity int    `json:"remainingQuantity" db:"remaining_quantity"`
	}
)

var SQL_GET_TOTAL_DISCOUNT = `
SELECT COUNT(*) AS TotalCount
FROM discounts
WHERE is_deleted = FALSE;
`

var SQL_GET_DISCOUNT = `
SELECT 
    id,
    code,
    amount,
    total_quantity,
    remaining_quantity
FROM discounts
WHERE is_deleted = FALSE
ORDER BY id ASC
`

var SQL_CREATE_DISCOUNT = `
INSERT INTO discounts (code, amount, total_quantity, remaining_quantity, created_at, updated_at, is_deleted)
VALUES (?, ?, ?, ?, NOW(), NOW(), FALSE);
`

var SQL_UPDATE_DISCOUNT = `
UPDATE discounts
SET code = ?, amount = ?, total_quantity = ?, remaining_quantity = ?, updated_at = NOW()
WHERE id = ? AND is_deleted = FALSE;
`

var SQL_SOFT_DELETE_DISCOUNT = `
UPDATE discounts
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = ? AND is_deleted = FALSE;
`
