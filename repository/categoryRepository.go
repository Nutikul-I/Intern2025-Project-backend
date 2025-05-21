package repository

import (
	"context"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

func GetCategoryRepository(mid string, page int, row int) ([]model.CategoryPlayload, error) {
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
		tsql = model.SQL_GET_CATEGORY + " WHERE CountChild > 0  LIMIT ? OFFSET ?"
	} else {
		tsql = model.SQL_GET_CATEGORY + " LIMIT ? OFFSET ?"
	}

	rows, err := conn.QueryContext(ctx, tsql, row, Offset)
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categoryData []model.CategoryPlayload
	err = scan.Rows(&categoryData, rows)
	if err != nil {
		log.Errorf("Error scanning rows: %v", err)
		return categoryData, err
	}

	if len(categoryData) == 0 {
		log.Infof("No categorys found for the given query.")
	}

	return categoryData, nil
}

func GetTotalCategoryRepository() (int, error) {
	conn := ConnectDB()
	ctx := context.Background()

	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return 0, err
	}

	row := conn.QueryRowContext(ctx, model.SQL_GET_TOTAL_CATEGORY)

	var total int
	err = row.Scan(&total)
	if err != nil {
		log.Errorf("Error scanning total count: %v", err)
		return 0, err
	}

	return total, nil
}

// Create category
func CreateCategoryRepository(body model.CreateCategory) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_CREATE_CATEGORY

	_, err := conn.ExecContext(ctx, tsql, body.Name)
	if err != nil {
		log.Errorf("Error executing insert: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "create category fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "create category success"}, nil
}

// Update category
func UpdateCategoryRepository(body model.UpdateCategory) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_UPDATE_CATEGORY
	_, err := conn.ExecContext(ctx, tsql,
		body.Name,
		body.Id,
	)
	if err != nil {
		log.Errorf("Error executing update: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "update category fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "update category success"}, nil
}

// Soft delete category
func DeleteCategoryRepository(categoryID int) (model.UpdateResponse, error) {
	conn := ConnectDB()
	ctx := context.Background()

	if err := conn.PingContext(ctx); err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql := model.SQL_SOFT_DELETE_CATEGORY
	_, err := conn.ExecContext(ctx, tsql, categoryID)
	if err != nil {
		log.Errorf("Error executing delete: %v", err)
		return model.UpdateResponse{StatusCode: 400, Message: "delete category fail"}, err
	}

	return model.UpdateResponse{StatusCode: 200, Message: "delete category success"}, nil
}
