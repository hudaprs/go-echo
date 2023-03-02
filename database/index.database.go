package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Credentials() (string, string, string, string, string) {
	databaseUsername := os.Getenv("DB_USERNAME")
	databasePassword := os.Getenv("DB_PASSWORD")
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")

	return databaseUsername, databasePassword, databaseHost, databasePort, databaseName
}

func MySQL() (*gorm.DB, error) {
	databaseUsername, databasePassword, databaseHost, databasePort, databaseName := Credentials()

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       databaseUsername + ":" + databasePassword + "@tcp(" + databaseHost + ":" + databasePort + ")" + "/" + databaseName + "?charset=utf8&parseTime=True&loc=Local", // Data Source Name
		DefaultStringSize:         256,                                                                                                                                                           // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                                                          // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                                                          // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                                                          // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                                                         // auto configure based on currently MySQL version
	}), &gorm.Config{})

	return db, err
}

func PostgreSQL() (*gorm.DB, error) {
	databaseUsername, databasePassword, _, _, databaseName := Credentials()

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=" + databaseUsername + " password=" + databasePassword + " dbname=" + databaseName + " port=5432 sslmode=disable TimeZone=Asia/Jakarta", // Data Source Name
		PreferSimpleProtocol: true,                                                                                                                                          // disables implicit prepared statement usage                                                                                                                                                    // auto configure based on currently MySQL version
	}), &gorm.Config{})

	return db, err
}

func DatabaseInit() (*gorm.DB, error) {
	databaseUsername, _, _, _, databaseName := Credentials()

	// MySQL
	// db, err = MySQL()

	// PostgreSQL
	db, err = PostgreSQL()

	if err != nil {
		panic("Database: Failed to connect to database")
	}

	fmt.Println("Database: connected to " + databaseName + " using " + databaseUsername)

	return db, err
}

func Connect() *gorm.DB {
	db := db.Debug()

	return db
}

func BeginTransaction(db *gorm.DB) (*gorm.DB, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	return tx, nil
}
