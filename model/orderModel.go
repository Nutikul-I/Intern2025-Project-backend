package model

/* ---------- รายการ (List) ---------- */
type OrderRow struct {
	OrderID      int64   `json:"OrderID"      db:"OrderID"`
	OrderCode    string  `json:"OrderCode"    db:"OrderCode"`
	CustomerName string  `json:"CustomerName" db:"CustomerName"`
	OrderDate    string  `json:"OrderDate"    db:"OrderDate"` // formatted TEXT
	ShippingFee  float64 `json:"ShippingFee"  db:"ShippingFee"`
	Discount     float64 `json:"Discount"     db:"Discount"`
	Total        float64 `json:"Total"        db:"Total"`
	Status       string  `json:"Status"       db:"StatusName"`
}

type OrderPagination struct {
	TotalPages int        `json:"TotalPages"`
	List       []OrderRow `json:"List"`
}

/* ---------- รายละเอียด (Detail) ---------- */
type OrderDetail struct {
	OrderID      int64            `json:"OrderID"`
	OrderCode    string           `json:"OrderCode"`
	Status       string           `json:"Status"`
	OrderDate    string           `json:"OrderDate"`
	CustomerName string           `json:"CustomerName"`
	Address      string           `json:"Address"`
	ShippingFee  float64          `json:"ShippingFee"`
	Discount     float64          `json:"Discount"`
	SubTotal     float64          `json:"SubTotal"`
	Total        float64          `json:"Total"`
	Items        []OrderItemBrief `json:"Items"`
}

type OrderItemBrief struct {
	ProductName string  `json:"ProductName" db:"ProductName"`
	Quantity    int     `json:"Qty"         db:"Qty"`
	Price       float64 `json:"Price"       db:"Price"`
	ImageURL    string  `json:"ImageURL"    db:"ImageURL"`
}

/* ---------- payload สำหรับ Update status ---------- */
type OrderUpdate struct {
	StatusID int64 `json:"StatusID"`
}

/* ---------- SQL ---------- */

/* --- List --- */
var SQL_GET_ORDER_LIST = `
SELECT
  o.id                                AS OrderID,
  LPAD(o.id,6,'0')                    AS OrderCode,
  CONCAT(u.first_name,' ',u.last_name) AS CustomerName,
  DATE_FORMAT(o.created_at,'%d/%m/%Y %H:%i') AS OrderDate,
  IFNULL(d.amount, 0)					AS Discount,
  0                                   AS ShippingFee,
  o.total_price                       AS Total,
  s.status_name                       AS StatusName
FROM orders o
JOIN customers     c ON c.id = o.customer_id
JOIN users         u ON u.id = c.user_id
JOIN order_statuses s ON s.id = o.order_status_id
LEFT JOIN discounts d ON d.id = o.discount_id
WHERE o.is_deleted = 0
ORDER BY o.id DESC
LIMIT ? OFFSET ?;
`

var SQL_GET_TOTAL_ORDERS = `SELECT COUNT(*) FROM orders WHERE is_deleted = 0;`

/* --- Detail (main row) --- */
var SQL_GET_ORDER_DETAIL = `
SELECT
  o.id            AS OrderID,
  LPAD(o.id,6,'0')AS OrderCode,
  s.status_name   AS Status,
  DATE_FORMAT(o.created_at,'%d/%m/%Y %H:%i') AS OrderDate,
  CONCAT(u.first_name,' ',u.last_name) AS CustomerName,
  CONCAT(a.address,' ',a.city,' ',a.district,' ',a.postal_code) AS Address,
  0               AS ShippingFee,
  IFNULL(d.amount,0)                   AS Discount,
  o.total_price                        AS SubTotal,
  (o.total_price - IFNULL(d.amount,0)) AS Total
FROM orders o
JOIN customers     c ON c.id = o.customer_id
JOIN users         u ON u.id = c.user_id
JOIN addresses     a ON a.id = o.address_id
JOIN order_statuses s ON s.id = o.order_status_id
LEFT JOIN discounts d ON d.id = o.discount_id
WHERE o.id = ? AND o.is_deleted = 0;
`

/* --- Items --- */
var SQL_GET_ORDER_ITEMS = `
SELECT
  p.name                              AS ProductName,
  i.quantity                          AS Qty,
  (i.quantity * p.price)              AS Price,
  COALESCE(
        (SELECT image_url
         FROM   product_images
         WHERE  product_id = p.id AND is_deleted = 0
         ORDER  BY id ASC
         LIMIT 1), ''
  )                                   AS ImageURL
FROM order_items i
JOIN products p ON p.id = i.product_id
WHERE i.order_id = ? AND i.is_deleted = 0;
`

/* --- Update status --- */
var SQL_UPDATE_ORDER_STATUS = `
UPDATE orders
SET order_status_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0;
`

/* --- Soft delete --- */
var SQL_SOFT_DEL_ORDER = `UPDATE orders       SET is_deleted = 1, deleted_at = NOW() WHERE id = ?`
var SQL_SOFT_DEL_ORDER_ITEMS = `UPDATE order_items  SET is_deleted = 1, deleted_at = NOW() WHERE order_id = ?`
