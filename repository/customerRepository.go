package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payso-internal-api/model"
	"strings"

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

func isUsernameAvailable(ctx context.Context, tx *sql.Tx, uname string) (bool, error) {
	var cnt int
	err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE user_name = ?`,
		uname,
	).Scan(&cnt)
	return cnt == 0, err
}

/* -------- CREATE -------- */
func CreateCustomer(ctx context.Context, in model.CustomerCreate) (int64, error) {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	// ตรวจ user_name ซ้ำ
	ok, err := isUsernameAvailable(ctx, tx, in.UserName)
	if err != nil || !ok {
		tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("username already taken")
	}

	pwHash, err := hash(in.Password)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	/* --- users --- */
	uRes, err := tx.ExecContext(ctx, model.SQL_INSERT_USER,
		in.FirstName, in.LastName, in.Email, in.Phone, in.UserName, pwHash)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	uid, _ := uRes.LastInsertId()

	/* --- customers --- */
	cRes, err := tx.ExecContext(ctx, model.SQL_INSERT_CUSTOMER, uid, in.NationalID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	cid, _ := cRes.LastInsertId()

	/* --- addresses --- */
	for _, a := range in.Addresses {
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_ADDRESS,
			cid, a.Name, a.Address, a.City, a.District, a.SubDistrict, a.PostalCode, a.Phone); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return cid, tx.Commit()
}

/* -------- UPDATE -------- */
func UpdateCustomer(ctx context.Context, cid int64, in model.CustomerCreate) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	/* 1. หา user_id */
	var uid int64
	if err := tx.QueryRowContext(ctx, model.SQL_SELECT_UID_BY_CID, cid).Scan(&uid); err != nil {
		tx.Rollback()
		return err
	}

	/* 2. ตรวจ user_name (ถ้าส่ง) */
	if uname := strings.TrimSpace(in.UserName); uname != "" {
		ok, err := isUsernameAvailable(ctx, tx, uname)
		if err != nil {
			tx.Rollback()
			return err
		}
		if !ok {
			tx.Rollback()
			return fmt.Errorf("username already taken")
		}
	}

	/* 3. อัปเดต users */
	if pw := strings.TrimSpace(in.Password); pw != "" {
		pwHash, err := hash(pw)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(ctx, model.SQL_UPDATE_USER_WITH_PW,
			in.FirstName, in.LastName, in.Email, in.Phone,
			in.UserName, pwHash, uid)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, model.SQL_UPDATE_USER_NO_PW,
			in.FirstName, in.LastName, in.Email, in.Phone,
			in.UserName, uid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	/* 4. อัปเดต national_id */
	if _, err := tx.ExecContext(ctx, model.SQL_UPDATE_CUSTOMER_NATID,
		in.NationalID, cid); err != nil {
		tx.Rollback()
		return err
	}

	/* 5. รีเฟรช addresses */
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ADDRESSES, cid); err != nil {
		tx.Rollback()
		return err
	}
	for _, a := range in.Addresses {
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_ADDRESS,
			cid, a.Name, a.Address, a.City, a.District, a.SubDistrict, a.PostalCode, a.Phone); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/* -------- DELETE -------- */
func DeleteCustomer(ctx context.Context, cid int64) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	/* ดึง user_id */
	var uid int64
	if err := tx.QueryRowContext(ctx, model.SQL_SELECT_UID_BY_CID, cid).Scan(&uid); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_USER, uid); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ADDRESSES, cid); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
