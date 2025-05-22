package repository

import (
	"context"
	"database/sql"
	"math"

	"payso-internal-api/model"

	"github.com/blockloop/scan"
)

/* ---------- List ---------- */
func GetOrderList(ctx context.Context, page, row int) (model.OrderPagination, error) {
	conn := ConnectDB()
	if row <= 0 {
		row = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * row

	rows, err := conn.QueryContext(ctx, model.SQL_GET_ORDER_LIST, row, offset)
	if err != nil {
		return model.OrderPagination{}, err
	}
	defer rows.Close()

	var list []model.OrderRow
	if err := scan.Rows(&list, rows); err != nil {
		return model.OrderPagination{}, err
	}

	var total int
	_ = conn.QueryRowContext(ctx, model.SQL_GET_TOTAL_ORDERS).Scan(&total)

	return model.OrderPagination{
		TotalPages: int(math.Ceil(float64(total) / float64(row))),
		List:       list,
	}, nil
}

/* ---------- Detail ---------- */
func GetOrderDetail(ctx context.Context, id int64) (model.OrderDetail, error) {
	conn := ConnectDB()

	var detail model.OrderDetail
	rows, err := conn.QueryContext(ctx, model.SQL_GET_ORDER_DETAIL, id)
	if err != nil {
		return detail, err
	}
	defer rows.Close()

	if err := scan.Row(&detail, rows); err != nil {
		return detail, err
	}

	// --- ดึงรายการ Items ---
	itemRows, _ := conn.QueryContext(ctx, model.SQL_GET_ORDER_ITEMS, id)
	defer itemRows.Close()
	for itemRows.Next() {
		var it model.OrderItemBrief
		_ = itemRows.Scan(&it.ProductName, &it.Quantity, &it.Price)
		detail.Items = append(detail.Items, it)
	}
	return detail, nil
}

/* ---------- Update Status ---------- */
func UpdateOrderStatus(ctx context.Context, id, statusID int64) error {
	conn := ConnectDB()
	_, err := conn.ExecContext(ctx, model.SQL_UPDATE_ORDER_STATUS, statusID, id)
	return err
}

/* ---------- Soft-delete ---------- */
func DeleteOrder(ctx context.Context, id int64) error {
	conn := ConnectDB()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ORDER, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, model.SQL_SOFT_DEL_ORDER_ITEMS, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
