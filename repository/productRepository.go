package repository

import (
	"context"
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

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("PingContext: %v", err)
		return model.ProductDetail{}, err
	}

	var detail model.ProductDetail

	rows, err := conn.QueryContext(ctx, model.SQL_GET_PRODUCT_DETAIL, pid)
	if err != nil {
		return detail, err
	}
	defer rows.Close()

	// ❌ อย่าเรียก rows.Next() ที่นี่
	if err := scan.Row(&detail, rows); err != nil {
		return detail, err // จะได้ sql.ErrNoRows ก็ต่อเมื่อไม่มีจริง ๆ
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

	return detail, nil
}
