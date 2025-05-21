package repository

import (
	"context"
	"database/sql"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	"golang.org/x/crypto/bcrypt"
)

func GetCustomerList(ctx context.Context, page, row int) ([]model.CustomerDetail, error) {
	conn := ConnectDB()
	if page <= 0 {
		page = 1
	}
	if row <= 0 {
		row = 20
	}
	offset := (page - 1) * row

	rows, err := conn.QueryContext(ctx, model.SQL_GET_CUSTOMER_LIST, row, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CustomerDetail
	if err := scan.Rows(&list, rows); err != nil {
		return nil, err
	}
	return list, nil
}

func GetCustomerDetail(ctx context.Context, id int64) (model.CustomerDetail, error) {
	conn := ConnectDB()
	var detail model.CustomerDetail

	/* ---------- main row ---------- */
	rows, err := conn.QueryContext(ctx, model.SQL_GET_CUSTOMER_DETAIL, id)
	if err != nil {
		return detail, err // query error (DB down ฯลฯ)
	}
	defer rows.Close()

	if err := scan.Row(&detail, rows); err != nil {
		return detail, err // รวมถึง sql.ErrNoRows
	}

	/* ---------- addresses ---------- */
	addrRows, _ := conn.QueryContext(ctx, model.SQL_GET_ADDRESS_BY_CUSTOMER, id)
	defer addrRows.Close()
	for addrRows.Next() {
		var a model.AddressPayload
		_ = addrRows.Scan(&a.ID, &a.Name, &a.Address, &a.City,
			&a.District, &a.SubDistrict, &a.PostalCode, &a.Phone)
		detail.Addresses = append(detail.Addresses, a)
	}
	return detail, nil
}

func hash(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

/* -------- CREATE -------- */
func CreateCustomer(ctx context.Context, c model.CustomerCreate) (int64, error) {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	pwHash, err := hash(c.Password)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO customers
           (first_name,last_name,email,phone,national_id,password_hash)
         VALUES (?,?,?,?,?,?)`,
		c.FirstName, c.LastName, c.Email, c.Phone, c.NationalID, pwHash)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	cid, _ := res.LastInsertId()

	for _, a := range c.Addresses {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO addresses
			  (customer_id,name,address,city,district,subdistrict,postal_code,phone)
			 VALUES (?,?,?,?,?,?,?,?)`,
			cid, a.Name, a.Address, a.City, a.District, a.SubDistrict, a.PostalCode, a.Phone,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return cid, tx.Commit()
}

/* -------- UPDATE -------- */
func UpdateCustomer(ctx context.Context, id int64, c model.CustomerCreate) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// ----- update main -----
	if c.Password != "" {
		pwHash, err := hash(c.Password)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE customers SET
                first_name = ?, last_name = ?, email = ?, phone = ?,
                national_id = ?, password_hash = ?, updated_at=CURRENT_TIMESTAMP
              WHERE id=? AND is_deleted=0`,
			c.FirstName, c.LastName, c.Email, c.Phone,
			c.NationalID, pwHash, id)
	} else {
		_, err = tx.ExecContext(ctx,
			`UPDATE customers SET
                first_name = ?, last_name = ?, email = ?, phone = ?,
                national_id = ?, updated_at=CURRENT_TIMESTAMP
              WHERE id=? AND is_deleted=0`,
			c.FirstName, c.LastName, c.Email, c.Phone,
			c.NationalID, id)
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	/* soft-delete old addresses */
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ADDRESSES, id); err != nil {
		tx.Rollback()
		return err
	}

	/* insert new addresses */
	for _, a := range c.Addresses {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO addresses
			  (customer_id,name,address,city,district,subdistrict,postal_code,phone)
			 VALUES (?,?,?,?,?,?,?,?)`,
			id, a.Name, a.Address, a.City, a.District, a.SubDistrict, a.PostalCode, a.Phone); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

/* -------- DELETE -------- */
func DeleteCustomer(ctx context.Context, id int64) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_CUSTOMER, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ADDRESSES, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
