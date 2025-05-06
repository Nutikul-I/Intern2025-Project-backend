package repository

import (
	"context"
	"database/sql"
	"payso-internal-api/model"

	"github.com/blockloop/scan"
	log "github.com/sirupsen/logrus"
)

func GetMerchantRepository(mid string, page int, row int) ([]model.MerchantPlayload, error) {
	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return nil, err
	}

	var Offset int = 0
	if page > 0 {
		Offset = (page - 1) * row
	}

	tsql := model.SQL_GET_MERCHANT
	if mid == "0" {
		tsql += " WHERE CountChild > 0 ORDER BY MerchantID DESC OFFSET @Offset ROWS FETCH NEXT @Row ROWS ONLY "
	}
	rows, err := conn.QueryContext(ctx, tsql,
		sql.Named("MerchantID", mid),
		sql.Named("Offset", Offset),
		sql.Named("Row", row))
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var MerchantData []model.MerchantPlayload
	err = scan.Rows(&MerchantData, rows)
	if err != nil {
		log.Errorf("Error scanning rows: %v", err)
		return MerchantData, err
	}

	return MerchantData, nil
}

func GetTotalMerchantRepository(mid string) (int, error) {

	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return 0, err
	}

	tsql := model.SQL_GET_TOTAL_MERCHANT
	rows, err := conn.QueryContext(ctx, tsql, sql.Named("MerchantID", mid))
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return 0, err
	}
	defer rows.Close()

	TotalCount := 0
	err = scan.Row(&TotalCount, rows)
	if err != nil {
		log.Errorf("Error scan row : %v", err)
		return 0, err
	}

	return TotalCount, nil
}

func CreateMerchantRepository(body model.CreateMerchantPayload) (model.UpdateResponse, error) {

	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql_check := model.SQL_CHECK_MERCHANT
	rows_check, err := conn.QueryContext(ctx, tsql_check,
		sql.Named("MasterMerchantID", body.MasterMerchantID),
		sql.Named("MerchantID", body.MerchantID))
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var MerchantData model.MasterMerchant
	err = scan.Row(&MerchantData, rows_check)
	if err != nil {
		tsql := model.SQL_CREATE_MERCHANT
		rows, err := conn.QueryContext(ctx, tsql,
			sql.Named("MasterMerchantID", body.MasterMerchantID),
			sql.Named("MerchantID", body.MerchantID))
		if err != nil {
			log.Errorf("Error executing query: %v", err)
			return model.UpdateResponse{}, err
		}
		defer rows.Close()

		return model.UpdateResponse{StatusCode: 200, Message: "created  merchant success"}, nil
	} else {
		return model.UpdateResponse{StatusCode: 400, Message: "created  merchant fail"}, nil
	}
}

func DeleteMerchantRepository(ReqMasterMerchantID string, ReqMerchantID string) (model.UpdateResponse, error) {

	conn := ConnectDB()
	ctx := context.Background()

	// Check if database is alive.
	err := conn.PingContext(ctx)
	if err != nil {
		log.Errorf("Error PingContext: %v", err)
		return model.UpdateResponse{}, err
	}

	tsql_check := model.SQL_CHECK_MERCHANT
	rows_check, err := conn.QueryContext(ctx, tsql_check,
		sql.Named("MasterMerchantID", ReqMasterMerchantID),
		sql.Named("MerchantID", ReqMerchantID))
	if err != nil {
		log.Errorf("Error executing query: %v", err)
		return model.UpdateResponse{}, err
	}
	defer rows_check.Close()

	var MerchantData model.MasterMerchant
	err = scan.Row(&MerchantData, rows_check)
	if err != nil {
		return model.UpdateResponse{StatusCode: 400, Message: "deleted  merchant fail"}, nil
	} else {
		tsql := model.SQL_DELETE_MERCHANT
		rows, err := conn.QueryContext(ctx, tsql,
			sql.Named("MasterMerchantID", ReqMasterMerchantID),
			sql.Named("MerchantID", ReqMerchantID))
		if err != nil {
			log.Errorf("Error executing query: %v", err)
			return model.UpdateResponse{}, err
		}
		defer rows.Close()
		return model.UpdateResponse{StatusCode: 200, Message: "deleted  merchant success"}, nil
	}
}
