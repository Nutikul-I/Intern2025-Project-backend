package model

type (
	UpdateResponse struct {
		StatusCode int    `json:"statuscode"`
		Message    string `json:"message"`
	}

	MerchantPlayload struct {
		MasterMerchantID *string `json:"MasterMerchantID"`
		MerchantID       *string `json:"MerchantID"`
		Company          *string `json:"Company"`
		Sale             *string `json:"Sale"`
		ActiveStatus     bool    `json:"ActiveStatus"`
		Indent           int     `json:"Indent"`
		Count            int     `json:"Count"`
		Child            []MerchantPlayload
	}

	MerchantPagination struct {
		TotalPages   int `json:"TotalPages"`
		MerchantList []MerchantPlayload
	}

	MerchantDetail struct {
		MerchantID      *string `json:"MerchantID"`
		MerchantCompany *string `json:"MerchantCompany"`
		NameTH          *string `json:"NameTH"`
		NameEN          *string `json:"NameEN"`
		SurnameTH       *string `json:"SurnameTH"`
		SurnameEN       *string `json:"SurnameEN"`
		AddressTH       *string `json:"AddressTH"`
		AddressEN       *string `json:"AddressEN"`
		Address2TH      *string `json:"Address2TH"`
		Address2EN      *string `json:"Address2EN"`
		CityTH          *string `json:"CityTH"`
		CityEN          *string `json:"CityEN"`
		Province        *string `json:"Province"`
		PostalcodeTH    *string `json:"PostalcodeTH"`
		CountryName     *string `json:"CountryName"`
		TelTH           *string `json:"TelTH"`
		Email           *string `json:"Email"`
	}

	CreateMerchantPayload struct {
		MasterMerchantID string `json:"MasterMerchantID"`
		MerchantID       string `json:"MerchantID"`
	}

	MasterMerchant struct {
		MasterMerchantID string `json:"MasterMerchantID"`
		MerchantID       string `json:"MerchantID"`
	}
)

var SQL_GET_MERCHANT = `
WITH MerchantData AS (
    SELECT
        mm.MasterMerchantID AS MasterMerchantID,
        m.MerchantID AS MerchantID,
        REPLACE(m.MerchantCompany, ',', '') AS Company,
        s.SaleName AS Sale,
        m.DemoStatus AS ActiveStatus,
        (SELECT COUNT(*) 
         FROM MerchantMaster AS a
         WHERE a.MasterMerchantID = m.MerchantID) AS CountChild
    FROM Merchant AS m
    LEFT JOIN MerchantMaster AS mm ON m.MerchantID = mm.MerchantID
    LEFT JOIN SaleData AS s ON m.MerchantID = s.MMID
    WHERE ((@MerchantID = '0') OR (m.MerchantID = @MerchantID))
)
SELECT *
FROM MerchantData`

var SQL_GET_TOTAL_MERCHANT = `
SELECT COUNT(*) AS TotalCount
FROM Merchant AS m
LEFT JOIN MerchantMaster AS mm ON m.MerchantID = mm.MerchantID
LEFT JOIN SaleData AS s ON m.MerchantID = s.MMID
WHERE ((@MerchantID = '0') OR (m.MerchantID = @MerchantID))`

var SQL_GET_CHECK_MERCHANT = `
SELECT 
m.MerchantID,
COALESCE(m.MerchantCompany, '') AS MerchantCompany,
COALESCE(m.NameTH, '') AS NameTH,
COALESCE(m.NameEN, '') AS NameEN,
COALESCE(m.SurnameTH, '') AS SurnameTH,
COALESCE(m.SurnameEN, '') AS SurnameEN,
COALESCE(m.AddressTH, '') AS AddressTH,
COALESCE(m.AddressEN, '') AS AddressEN,
COALESCE(m.Address2TH, '') AS Address2TH,
COALESCE(m.Address2EN, '') AS Address2EN,
COALESCE(m.CityTH, '') AS CityTH,
COALESCE(m.CityEN, '') AS CityEN,
COALESCE(p.NameTH, '') AS Province,
COALESCE(m.PostalcodeTH, '') AS PostalcodeTH,
COALESCE(c.Name, '') AS CountryName,
COALESCE(m.TelTH, '') AS TelTH,
COALESCE(m.Email, '') AS Email
FROM Merchant AS m 
LEFT JOIN Province AS p ON m.ProvinceID = p.ProvinceID
LEFT JOIN Country AS c ON m.CountryCodeTH = c.Code
WHERE m.MerchantID = @MerchantID`

var SQL_CHECK_MERCHANT = `
SELECT 
m.MasterMerchantID,
m.MerchantID
FROM MerchantMaster AS m 
WHERE m.MerchantID = @MerchantID`

var SQL_CREATE_MERCHANT = `INSERT INTO MerchantMaster (MasterMerchantID, MerchantID) VALUES (@MasterMerchantID, @MerchantID)`

var SQL_DELETE_MERCHANT = `DELETE FROM MerchantMaster WHERE MasterMerchantID = @MasterMerchantID AND MerchantID = @MerchantID;`

var SQL_COUNT_MERCHANT = `SELECT COUNT(*) from Merchant m 
WHERE m.MerchantID = @MerchantID`
