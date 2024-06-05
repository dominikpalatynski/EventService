package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn string = "host=localhost  user=postgres_user password=s3cr3t dbname=postgres_db port=5432 sslmode=disable"

type StorePostgres struct {
	db *gorm.DB
}

func NewPostgresStore() (*StorePostgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &StorePostgres {
		db: db}, nil
}