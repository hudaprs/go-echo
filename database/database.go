package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func DatabaseInit() (*gorm.DB, error) {
	databaseUsername := os.Getenv("DB_USERNAME")
	databasePassword := os.Getenv("DB_PASSWORD")
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       databaseUsername + ":" + databasePassword + "@tcp(" + databaseHost + ":" + databasePort + ")" + "/" + databaseName + "?charset=utf8&parseTime=True&loc=Local", // Data Source Name
		DefaultStringSize:         256,                                                                                                                                                           // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                                                          // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                                                          // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                                                          // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                                                         // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		panic("Database: Failed to connect to database")
	}

	fmt.Println("Database: connected to " + databaseName + " using " + databaseUsername)

	return db, err
}

func DatabaseConnection() *gorm.DB {
	return db
}
