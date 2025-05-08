package db

import (
	"database/sql"

	"github.com/DjentBoiiii/marketplace/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("mysql", config.DB_USER+":"+config.DB_PASSWORD+"@tcp("+config.DB_HOST+":3306)/"+config.DB_NAME)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		panic(err)
	}

}
