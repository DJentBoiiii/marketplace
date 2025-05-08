package config

import (
	"os"

	"github.com/DjentBoiiii/marketplace/internal/db"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
)

var (
	_           = godotenv.Load("/marketplace/.env")
	DB_USER     = os.Getenv("MYSQL_USER")
	DB_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	DB_NAME     = os.Getenv("MYSQL_DATABASE")
	JWT_SECRET  = os.Getenv("JWT_SECRET")
	SHA_SECRET  = os.Getenv("SHA_SECRET")
	DB_HOST     = os.Getenv("DB_HOST")
	DB          = db.DB
)

var (
	MinioEndpoint  = os.Getenv("MINIO_ENDPOINT")
	MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	MinioClient    *minio.Client
)
