package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Replace with your own connection parameters
var server = "x"
var port = 1433
var user = "x"
var password = "x"
var database = "x"
var DB *sql.DB

func Init() {
	server = viper.GetString("DB_SERVER")
	port = 1433
	user = viper.GetString("DB_USER")
	password = viper.GetString("DB_PASS")
	database = viper.GetString("DB_INST")
	// Create connection string
	var err error

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

	// Create connection pool
	DB, err = sql.Open("sqlserver", connString)

	if err != nil {
		log.Error("**** Error creating connection pool: " + err.Error())
	}
	log.Debug("==-- Connected! --==")

}

func ConnectDB() *sql.DB {
	var err error

	if DB == nil {
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

		DB, err = sql.Open("sqlserver", connString)
		if err != nil {
			log.Error("**** Error creating connection pool: " + err.Error())
		}
	}

	log.Debug("==-- Connected! --==")
	return DB
}
