package database

import (
	"database/sql"
	"fmt"
	"os"

	_"github.com/go-sql-driver/mysql"
)

// Expose the database connection globally
var DB *sql.DB

func Init() (*sql.DB, error)  {
	// Initialize the database
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", dbuser + ":" + dbpass + "@tcp(localhost:3306)/" + dbname)
	if err != nil {
		fmt.Println(err)
	}
	DB = db

	return db, err
}
