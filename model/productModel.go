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
		ProductID    *int64      `json:"ProductID"`
		Name         *string     `json:"Name"`
		Description  *string     `json:"Description"`
		Price        *float64    `json:"Price"`
		UnitName     *string     `json:"UnitName"`
		CategoryName *string     `json:"CategoryName"`
		Images       []string    `json:"Images"`
		Colors       []string    `json:"Colors"`
		Attributes   []Attribute `json:"Attributes"`
	}

	Attribute struct {
		Name  *string `json:"Name"  db:"AttributeName"`
		Value *string `json:"Value" db:"AttributeValue"`
	}
)

/* ----------------- CHILD PAYLOAD ----------------- */

// สีที่เลือก
type AttributePayload struct {
	Name  string `json:"Name"`  // หัวข้อ   เช่น "ขนาดหน้าจอ"
	Value string `json:"Value"` // รายละเอียด เช่น "6.7\" OLED"
}

// สำหรับสี
type ColorPayload struct {
	Name string `json:"Name"` // เช่น "แดง"
	Hex  string `json:"Hex"`  // เช่น "#FF0000"
}

/* ----------------- MAIN PAYLOAD ----------------- */

// POST / PUT (สร้าง-แก้ไขสินค้า)
type ProductCreate struct {
	Name        string             `json:"Name"`
	Description string             `json:"Description"`
	Price       float64            `json:"Price"`
	UnitID      int64              `json:"UnitID"`
	CategoryID  int64              `json:"CategoryID"`
	Images      []string           `json:"Images"`
	Colors      []ColorPayload     `json:"Colors"`
	Attributes  []AttributePayload `json:"Attributes"`
}

/* UPDATE product main */
var SQL_UPDATE_PRODUCT = `
UPDATE products
SET    name        = ?, 
       description = ?, 
       price       = ?, 
       unit_id     = ?, 
       category_id = ?, 
       updated_at  = CURRENT_TIMESTAMP
WHERE  id = ? AND is_deleted = 0
`

/* SOFT-DELETE children ก่อน insert ใหม่ */
var SQL_SOFT_DEL_IMAGES = `UPDATE product_images      SET is_deleted = 1, deleted_at = NOW() WHERE product_id = ? AND is_deleted = 0`
var SQL_SOFT_DEL_COLORS = `UPDATE product_colors      SET is_deleted = 1, deleted_at = NOW() WHERE product_id = ? AND is_deleted = 0`
var SQL_SOFT_DEL_ATTR = `UPDATE product_attributes  SET is_deleted = 1, deleted_at = NOW() WHERE product_id = ? AND is_deleted = 0`

/* SOFT-DELETE main + children (Delete API) */
var SQL_SOFT_DEL_PRODUCT = `UPDATE products SET is_deleted = 1, deleted_at = NOW() WHERE id = ? AND is_deleted = 0`

/* ---------------- SQL ---------------- */

/* ---------- LIST ---------- */
var SQL_GET_PRODUCT = `
SELECT
  p.id           AS ProductID,
  p.name         AS Name,
  p.description  AS Description,
  p.price        AS Price,
  c.name         AS CategoryName,
  0              AS StockQty
FROM products p
LEFT JOIN categories c ON c.id = p.category_id AND c.is_deleted = 0
WHERE (? = 0 OR p.id = ?)
  AND p.is_deleted = 0
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
LEFT JOIN units      u ON u.id = p.unit_id     AND u.is_deleted = 0
LEFT JOIN categories c ON c.id = p.category_id AND c.is_deleted = 0
WHERE p.id = ?          -- id ที่ต้องการ
  AND p.is_deleted = 0; -- กรอง soft-delete
`

/* ---------- CHILDREN ---------- */
var SQL_GET_PRODUCT_IMAGES = `
SELECT image_url
FROM   product_images
WHERE  product_id = ? AND is_deleted = 0
ORDER  BY id
`
var SQL_GET_PRODUCT_COLORS = `
SELECT color_name
FROM   product_colors
WHERE  product_id = ? AND is_deleted = 0
ORDER  BY id
`
var SQL_GET_PRODUCT_ATTRIB = `
SELECT attribute_name  AS AttributeName,
       attribute_value AS AttributeValue
FROM   product_attributes
WHERE  product_id = ? AND is_deleted = 0
ORDER  BY id
`

/* ---------- INSERT ---------- */
var SQL_INSERT_PRODUCT = `
INSERT INTO products
  (name, description, price, unit_id, category_id)
VALUES (?,?,?,?,?)`

var SQL_INSERT_PRODUCT_IMAGE = `INSERT INTO product_images  (product_id, image_url)                       VALUES (?,?)`
var SQL_INSERT_PRODUCT_COLOR = `INSERT INTO product_colors  (product_id, color_name, color_code)           VALUES (?,?,?)`
var SQL_INSERT_PRODUCT_ATTR = `INSERT INTO product_attributes (product_id, attribute_name, attribute_value) VALUES (?,?,?)`
