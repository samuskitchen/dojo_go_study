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
	db, err := getConnection()
	if err != nil {
		log.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		log.Println("We are connected to the database")
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = makeMigration(db)
	if err != nil {
		log.Fatal("This is the error:", err)
	}

	data = &Data{
		DB: db,
	}
}

func getConnection() (*sql.DB, error) {

	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	return sql.Open(DbDriver, uri)
}

func makeMigration(db *sql.DB) error {
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