package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Asia%%2FBangkok",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_SERVER"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_INST"),
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("เปิด connection ไม่ได้: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Ping DB ล้มเหลว: %v", err)
	}
	log.Info("🎉 MySQL connected")
}

func ConnectDB() *sql.DB { return DB }
