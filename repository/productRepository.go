package repository

import (
	"context"
	"database/sql"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

/* ---------- List ---------- */

func GetProductRepository(pid int64, page, row int) ([]model.ProductPayload, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("PingContext: %v", err)
		return nil, err
	}

	offset := 0
	if page > 0 {
		offset = (page - 1) * row
	}

	// SQL_GET_PRODUCT มี ORDER BY + LIMIT ? OFFSET ? อยู่แล้ว
	rows, err := conn.QueryContext(ctx, model.SQL_GET_PRODUCT,
		pid,    // ? #1  (WHERE (? = 0 OR id = ?))
		pid,    // ? #2
		row,    // ? #3  LIMIT
		offset, // ? #4  OFFSET
	)
	if err != nil {
		log.Errorf("exec list: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []model.ProductPayload
	if err := scan.Rows(&list, rows); err != nil {
		log.Errorf("scan list: %v", err)
		return nil, err
	}
	return list, nil
}

/* ---------- Total count ---------- */

func GetTotalProductRepository(pid int64) (int, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("PingContext: %v", err)
		return 0, err
	}

	row := conn.QueryRowContext(ctx, model.SQL_GET_TOTAL_PRODUCT,
		pid, // ? #1
		pid, // ? #2
	)

	var total int
	if err := row.Scan(&total); err != nil {
		log.Errorf("scan total: %v", err)
		return 0, err
	}
	return total, nil
}

/* ---------- Detail ---------- */

func GetProductDetailRepository(pid int64) (model.ProductDetail, error) {
	conn := ConnectDB()
	ctx := context.Background()

	var detail model.ProductDetail

	/* --- main detail --- */
	rows, err := conn.QueryContext(ctx, model.SQL_GET_PRODUCT_DETAIL, pid)
	if err != nil {
		return detail, err
	}
	defer rows.Close()
	if err := scan.Row(&detail, rows); err != nil {
		return detail, err // sql.ErrNoRows ถ้าไม่พบ
	}

	/* --- images --- */
	imgRows, _ := conn.QueryContext(ctx, model.SQL_GET_PRODUCT_IMAGES, pid)
	defer imgRows.Close()
	for imgRows.Next() {
		var url string
		_ = imgRows.Scan(&url)
		detail.Images = append(detail.Images, url)
	}

	/* --- colors --- */
	colRows, _ := conn.QueryContext(ctx, model.SQL_GET_PRODUCT_COLORS, pid)
	defer colRows.Close()
	for colRows.Next() {
		var c string
		_ = colRows.Scan(&c)
		detail.Colors = append(detail.Colors, c)
	}

	/* --- attributes (NEW) --- */
	attrRows, _ := conn.QueryContext(ctx, model.SQL_GET_PRODUCT_ATTRIB, pid)
	defer attrRows.Close()
	for attrRows.Next() {
		var a model.Attribute
		_ = attrRows.Scan(&a.Name, &a.Value)
		detail.Attributes = append(detail.Attributes, a)
	}

	return detail, nil
}

/* ---------- Create ---------- */

func CreateProductRepository(ctx context.Context, p model.ProductCreate) (int64, error) {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	/* --- main product --- */
	res, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT,
		p.Name, p.Description, p.Price, p.UnitID, p.CategoryID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	pid, _ := res.LastInsertId()

	/* --- images --- */
	for _, url := range p.Images {
		if url == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_IMAGE, pid, url); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	/* --- colors --- */
	for _, c := range p.Colors {
		if c.Name == "" && c.Hex == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_COLOR, pid, c.Name, c.Hex); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	/* --- attributes --- */
	for _, a := range p.Attributes {
		if a.Name == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_ATTR, pid, a.Name, a.Value); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	log.Infof("CreateProductRepository: new product_id=%d", pid)
	return pid, nil
}

/* ---------- Update ---------- */

func UpdateProductRepository(ctx context.Context, pid int64, p model.ProductCreate) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	/* 1. update main */
	if _, err := tx.ExecContext(ctx, model.SQL_UPDATE_PRODUCT,
		p.Name, p.Description, p.Price, p.UnitID, p.CategoryID, pid); err != nil {
		tx.Rollback()
		return err
	}

	/* 2. soft-delete children */
	for _, q := range []string{
		model.SQL_SOFT_DEL_IMAGES,
		model.SQL_SOFT_DEL_COLORS,
		model.SQL_SOFT_DEL_ATTR,
	} {
		if _, err := tx.ExecContext(ctx, q, pid); err != nil {
			tx.Rollback()
			return err
		}
	}

	/* 3. re-insert children */
	for _, url := range p.Images {
		if url == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_IMAGE, pid, url); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, c := range p.Colors {
		if c.Name == "" && c.Hex == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_COLOR, pid, c.Name, c.Hex); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, a := range p.Attributes {
		if a.Name == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, model.SQL_INSERT_PRODUCT_ATTR, pid, a.Name, a.Value); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/* ---------- DELETE (soft-delete) ---------- */
func DeleteProductRepository(ctx context.Context, pid int64) error {
	conn, err := ConnectDB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	tx := conn

	// 1. soft-delete main
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_PRODUCT, pid); err != nil {
		tx.Rollback()
		return err
	}

	// 2. soft-delete children
	for _, q := range []string{
		model.SQL_SOFT_DEL_IMAGES,
		model.SQL_SOFT_DEL_COLORS,
		model.SQL_SOFT_DEL_ATTR,
	} {
		if _, err := tx.ExecContext(ctx, q, pid); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
