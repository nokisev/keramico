package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	var err error

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	dbUser, _ := os.LookupEnv("DB_USER")
	dbPassword, _ := os.LookupEnv("DB_PASSWORD")
	dbHost, _ := os.LookupEnv("DB_HOST")
	dbPort, _ := os.LookupEnv("DB_PORT")
	dbName, _ := os.LookupEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MySQL database!")
}
