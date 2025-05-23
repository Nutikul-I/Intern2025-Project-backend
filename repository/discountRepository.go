package repository

import (
	"context"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

func GetDiscountRepository(mid string, page int, row int) ([]model.DiscountPlayload, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return nil, err
	}

	Offset := (page - 1) * row
	var tsql string

	if mid == "0" {
		tsql = model.SQL_GET_DISCOUNT + " WHERE CountChild > 0  LIMIT ? OFFSET ?"
	} else {
		tsql = model.SQL_GET_DISCOUNT + " LIMIT ? OFFSET ?"
	}

	rows, err := conn.QueryContext(ctx, tsql, row, Offset)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var discountData []model.DiscountPlayload
	err = scan.Rows(&discountData, rows)
	if err != nil {
		log.Errorf("Error scanning rows: %v", err)
		return discountData, err
	}

	if len(discountData) == 0 {
		log.Infof("No discounts found for the given query.")
	}

	return discountData, nil
}

func GetTotalDiscountRepository() (int, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return 0, err
	}

	row := conn.QueryRowContext(ctx, model.SQL_GET_TOTAL_DISCOUNT)

	var total int
	err = row.Scan(&total)
	if err != nil {
		log.Errorf("Error scanning total count: %v", err)
		return 0, err
	}

	return total, nil
}

// Create discount
func CreateDiscountRepository(body model.CreateDiscount) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_CREATE_DISCOUNT

	_, err := conn.ExecContext(ctx, tsql,
		body.Code,
		body.Amount,
		body.TotalQuantity,
		body.RemainingQuantity,
	)
	if err != nil {
		log.Errorf("Error executing insert: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "create discount fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "create discount success"}, nil
}

// Update discount
func UpdateDiscountRepository(body model.UpdateDiscount) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_UPDATE_DISCOUNT
	_, err := conn.ExecContext(ctx, tsql,
		body.Code,
		body.Amount,
		body.TotalQuantity,
		body.Id,
	)
	if err != nil {
		log.Errorf("Error executing update: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "update discount fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "update discount success"}, nil
}

// Soft delete discount
func DeleteDiscountRepository(discountID int) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_SOFT_DELETE_DISCOUNT
	_, err := conn.ExecContext(ctx, tsql, discountID)
	if err != nil {
		log.Errorf("Error executing delete: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "delete discount fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "delete discount success"}, nil
}
