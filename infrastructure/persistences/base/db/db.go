package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pm/domain/entity"
)

func NewDBConnection(dsn string) (*gorm.DB, error) {
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.Migrator().AutoMigrate(entity.Category{}, entity.Product{}); err != nil {
		fmt.Printf("error migrating entity: %v", err)

	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Category{},
		&entity.Product{},
		//&entity.User{},
	)
}

func SetupDatabase(dsn string) (*gorm.DB, error) {
	db, err := NewDBConnection(dsn)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func GetDSN(username, password, domain, port, dbName string) string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full",
		username,
		password,
		domain,
		port,
		dbName,
	)
	return dsn
}