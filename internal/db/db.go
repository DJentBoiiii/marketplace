package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
	SHA_SECRET  = os.Getenv("SHA_SECRET")
	DB_HOST     = os.Getenv("DB_HOST")
	DB          *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":3306)/"+DB_NAME)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		panic(err)
	}

}
