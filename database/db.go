package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {

	db, err := connectingDB()

	if err != nil {
		log.Error().Err(err).Msg("Not able to connect the DB, please Look at the db.go/InitDB")
		return err
	}
	DB = db
	return nil
}

func connectingDB() (*gorm.DB, error) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("Not able to connect the database")
		return nil, err
	}

	return db, nil
}
