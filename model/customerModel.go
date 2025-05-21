package model

/* ============ CHILD ============ */
type AddressPayload struct {
	ID          *int64 `json:"ID,omitempty"` // ใช้ตอน update
	Name        string `json:"Name"`
	Address     string `json:"Address"`
	City        string `json:"City"`
	District    string `json:"District"`
	SubDistrict string `json:"SubDistrict"`
	PostalCode  string `json:"PostalCode"`
	Phone       string `json:"Phone"`
}

/* ============ CREATE / UPDATE PAYLOAD ============ */
type CustomerCreate struct {
	FirstName  string           `json:"FirstName"`
	LastName   string           `json:"LastName"`
	Email      string           `json:"Email"`
	Phone      string           `json:"Phone"`
	NationalID string           `json:"NationalID"`
	UserName   string           `json:"UserName"`
	Password   string           `json:"Password,omitempty"`
	Addresses  []AddressPayload `json:"Addresses"`
}

/* ============ RESPONSE DTO ============ */
type CustomerDetail struct {
	ID         *int64           `json:"ID"`
	FirstName  *string          `json:"FirstName"`
	LastName   *string          `json:"LastName"`
	Email      *string          `json:"Email"`
	Phone      *string          `json:"Phone"`
	NationalID *string          `json:"NationalID"`
	Addresses  []AddressPayload `json:"Addresses"`
}

/* ============ SQL ============ */

/* ---------- LIST ---------- */
var SQL_GET_CUSTOMER_LIST = `
SELECT
  c.id          AS ID,
  u.first_name  AS FirstName,
  u.last_name   AS LastName,
  u.email       AS Email,
  u.phone       AS Phone,
  c.national_id AS NationalID
FROM customers c
JOIN users u ON u.id = c.user_id AND u.is_deleted = 0
WHERE u.is_deleted = 0
ORDER BY c.id DESC
LIMIT ? OFFSET ?;
`

var SQL_GET_TOTAL_CUSTOMER = `
SELECT COUNT(*) FROM customers WHERE is_deleted = 0;
`

/* ---------- DETAIL ---------- */
var SQL_GET_CUSTOMER_DETAIL = `
SELECT
  c.id          AS ID,
  u.first_name  AS FirstName,
  u.last_name   AS LastName,
  u.email       AS Email,
  u.phone       AS Phone,
  c.national_id AS NationalID
FROM customers c
JOIN users u ON u.id = c.user_id AND u.is_deleted = 0
WHERE c.id = ? AND u.is_deleted = 0;
`

var SQL_GET_ADDRESS_BY_CUSTOMER = `
SELECT id            AS ID,
       name          AS Name,
       address       AS Address,
       city          AS City,
       district      AS District,
       subdistrict   AS SubDistrict,
       postal_code   AS PostalCode,
       phone         AS Phone
FROM addresses
WHERE customer_id = ? AND is_deleted = 0
ORDER BY id;
`

/* ---------- SOFT-DELETE ---------- */
var SQL_SOFT_DEL_USER = `UPDATE users SET is_deleted = 1, deleted_at = NOW() WHERE id = ?`
var SQL_SOFT_DEL_ADDRESSES = `UPDATE addresses SET is_deleted = 1, deleted_at = NOW() WHERE customer_id = ?`

/* ---------- INSERT ---------- */
var SQL_INSERT_USER = `
INSERT INTO users
  (first_name, last_name, email, phone, user_name, password)
VALUES (?,?,?,?,?,?)`

var SQL_INSERT_CUSTOMER = `
INSERT INTO customers (user_id, national_id)
VALUES (?,?)`

var SQL_INSERT_ADDRESS = `
INSERT INTO addresses
  (customer_id, name, address, city, district, subdistrict, postal_code, phone)
VALUES (?,?,?,?,?,?,?,?)`

/* ---------- SELECT ---------- */
var SQL_SELECT_UID_BY_CID = `SELECT user_id FROM customers WHERE id = ?`

/* ---------- UPDATE ---------- */
var SQL_UPDATE_USER_WITH_PW = `
UPDATE users SET
  first_name = ?, last_name = ?, email = ?, phone = ?,
  user_name = ?, password = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0`

var SQL_UPDATE_USER_NO_PW = `
UPDATE users SET
  first_name = ?, last_name = ?, email = ?, phone = ?,
  user_name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0`

var SQL_UPDATE_CUSTOMER_NATID = `
UPDATE customers SET national_id = ?
WHERE id = ?`
