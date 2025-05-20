package model

type (
	// สำหรับ list
	ProductPayload struct {
		ProductID    *int64   `json:"ProductID"    db:"ProductID"`
		Name         *string  `json:"Name"         db:"Name"`
		Description  *string  `json:"Description"  db:"Description"`
		Price        *float64 `json:"Price"        db:"Price"`
		CategoryName *string  `json:"CategoryName" db:"CategoryName"`
		StockQty     *int     `json:"StockQty"     db:"StockQty"`
	}

	ProductPagination struct {
		TotalPages  int              `json:"TotalPages"`
		ProductList []ProductPayload `json:"ProductList"`
	}

	// สำหรับ detail
	ProductDetail struct {
		ProductID    *int64   `json:"ProductID"`
		Name         *string  `json:"Name"`
		Description  *string  `json:"Description"`
		Price        *float64 `json:"Price"`
		UnitName     *string  `json:"UnitName"`
		CategoryName *string  `json:"CategoryName"`
		Images       []string `json:"Images"`
		Colors       []string `json:"Colors"`
	}
)

/* ---------------- SQL ---------------- */

/* สำหรับ list */
var SQL_GET_PRODUCT = `
SELECT
  p.id           AS ProductID,
  p.name         AS Name,
  p.description  AS Description,
  p.price        AS Price,
  c.name         AS CategoryName,
  0              AS StockQty        -- ต้องมี alias
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE (? = 0 OR p.id = ?)
ORDER BY p.id DESC
LIMIT ? OFFSET ?;
`

/* นับทั้งหมด */
var SQL_GET_TOTAL_PRODUCT = `
SELECT COUNT(*)
FROM products
WHERE (? = 0 OR id = ?)
`

/* ดึงรายละเอียด */
var SQL_GET_PRODUCT_DETAIL = `
SELECT
  p.id           AS ProductID,
  p.name         AS Name,
  p.description  AS Description,
  p.price        AS Price,
  u.name         AS UnitName,
  c.name         AS CategoryName
FROM products p
LEFT JOIN units      u ON p.unit_id     = u.id
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.id = ?
`

var SQL_GET_PRODUCT_IMAGES = `SELECT image_url  FROM product_images  WHERE product_id = ?`
var SQL_GET_PRODUCT_COLORS = `SELECT color_name FROM product_colors WHERE product_id = ?`
