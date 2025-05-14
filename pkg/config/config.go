package config

import "os"

var (
	Host     = os.Getenv("DB_HOST")
	Port     = os.Getenv("DB_PORT")
	User     = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
	Dbname   = os.Getenv("DB_NAME")
)


