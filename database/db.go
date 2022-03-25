package database

import (
	"database/sql"
	"fmt"

	"fake.com/GoRPCApi/errhelp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := GetDatabase()
	if errhelp.ErrorExists(err) {
		errhelp.ThrowConnectionError(err)
	}
	return db
}

func GetDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", DB_HOST, USER, PASSWORD, DB_NAME, PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func CloseDB(db *gorm.DB) {
	sqlDb, err := db.DB()
	if errhelp.ErrorExists(err) {
		errhelp.ThrowCloseError(err)
	} else {
		closeSqlDB(sqlDb)
	}
}

func closeSqlDB(sqlDb *sql.DB) {
	sqlDb.Close()
}
