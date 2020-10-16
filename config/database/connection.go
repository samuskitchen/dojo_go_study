package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	// registering database driver
	_ "github.com/lib/pq"
)

var (
	data *Data
	once sync.Once
)

type Data struct {
	DB *sql.DB
}

// New returns a new instance of Data with the database connection ready.
func New() *Data {
	once.Do(initDB)
	return data
}

func initDB() {
	db, err := GetConnection()
	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	err = MakeMigration(db)
	if err != nil {
		log.Fatal("This is the error:", err)
	}

	data = &Data{
		DB: db,
	}
}

func GetConnection() (*sql.DB, error) {

	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	return sql.Open(DbDriver, uri)
}

func MakeMigration(db *sql.DB) error {
	file, err := ioutil.ReadFile("./config/database/models.sql")
	if err != nil {
		return err
	}

	rows, err := db.Query(string(file))
	if err != nil {
		return err
	}

	return rows.Close()
}

// Close closes the resources used by data.
func Close() error {
	if data == nil {
		return nil
	}

	return data.DB.Close()
}