package repository

import (
	"context"
	"database/sql"
)

// Queryer ถูกใช้ให้ generic ข้าม DB / mock unit-test ได้
// *sql.DB และ *sql.Tx จะ implement ฟังก์ชันเหล่านี้อยู่แล้ว
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
