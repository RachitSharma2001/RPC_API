package main

import (
	"database/sql"
	"fmt"

	errHelp "fake.com/GoRPCApi/ErrHelp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var err error
	db, err := GetDatabase()
	if errHelp.ErrorExists(err) {
		errHelp.ThrowConnectionError(err)
	}
	return db
}

func GetDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", DB_HOST, USER, PASSWORD, DB_NAME, PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func CloseDB() {
	sqlDb, err := db.DB()
	if errHelp.ErrorExists(err) {
		errHelp.ThrowCloseError(err)
	} else {
		closeSqlDB(sqlDb)
	}
}

func closeSqlDB(sqlDb *sql.DB) {
	sqlDb.Close()
}
