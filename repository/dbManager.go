package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Warn("Could not load .env file, using system env variables")
	}

	server, portStr, user, password, database := SetDataBaseEnv()

	port := 3306
	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Error("Invalid port number in environment variable: " + err.Error())
		port = 3306
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Error("**** Error creating connection pool: " + err.Error())
		return
	}
	log.Debug("==-- Connected! --==")
}

func ConnectDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Warn("Could not load .env file, using system env variables")
	}

	server, portStr, user, password, database := SetDataBaseEnv()

	port := 3306
	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Error("Invalid port number in environment variable: " + err.Error())
		port = 3306
	}

	if DB == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)

		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Error("**** Error creating connection pool: " + err.Error())
		}
	}

	log.Debug("==-- Connected! --==")
	return DB
}

func SetDataBaseEnv() (string, string, string, string, string) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	return os.Getenv("DB_SERVER"),
		os.Getenv("SERVER_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_INST")
}
