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

func NewTest() *Data {
	once.Do(initDBTest)
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

func initDBTest() {

	db, err := getConnectionTest()
	if err != nil {
		fmt.Println("Cannot connect to database test")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database test")
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = makeMigrationTest(db)
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

func getConnectionTest() (*sql.DB, error) {
	DbHost := os.Getenv("TestDbHost")
	DbDriver := os.Getenv("TestDbDriver")
	DbUser := os.Getenv("TestDbUser")
	DbPassword := os.Getenv("TestDbPassword")
	DbName := os.Getenv("TestDbName")
	DbPort := os.Getenv("TestDbPort")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	return sql.Open(DbDriver, uri)
}

// makeMigration creates all the tables in the database
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

// makeMigrationTest creates all the tables in the database
func makeMigrationTest(db *sql.DB) error {
	b, err := ioutil.ReadFile("../../config/database/models.sql")
	if err != nil {
		return err
	}

	rows, err := db.Query(string(b))
	if err != nil {
		return err
	}

	return rows.Close()
}