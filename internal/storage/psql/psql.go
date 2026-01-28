package psql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(storagePath string) (*Storage, error) {

	psqlDb, err := gorm.Open(postgres.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: psqlDb}, nil
}
