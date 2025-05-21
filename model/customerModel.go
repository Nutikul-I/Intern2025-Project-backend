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

var SQL_GET_CUSTOMER_LIST = `
SELECT id,
       first_name     AS FirstName,
       last_name      AS LastName,
       email          AS Email,
       phone          AS Phone,
       national_id    AS NationalID
FROM customers
WHERE is_deleted = 0
ORDER BY id DESC
LIMIT ? OFFSET ?;
`

var SQL_GET_TOTAL_CUSTOMER = `
SELECT COUNT(*) FROM customers WHERE is_deleted = 0;
`

var SQL_GET_CUSTOMER_DETAIL = `
SELECT id,
       first_name     AS FirstName,
       last_name      AS LastName,
       email          AS Email,
       phone          AS Phone,
       national_id    AS NationalID
FROM customers
WHERE id = ? AND is_deleted = 0;
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

/* --- soft-delete --- */
var SQL_SOFT_DEL_CUSTOMER = `UPDATE customers  SET is_deleted = 1, deleted_at = NOW() WHERE id = ?`
var SQL_SOFT_DEL_ADDRESSES = `UPDATE addresses  SET is_deleted = 1, deleted_at = NOW() WHERE customer_id = ?`
