package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"leoho.io/it16th-webauthn-rp-server/config"
)

type context struct {
	mu sync.Mutex
	db *gorm.DB
}

var Context *context

func Connect() {
	dbConfig := config.GetDatabaseConfiguration()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DatabaseName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	migrator := db.Migrator()
	migrator.HasTable(&User{})

	Context = &context{db: db}
}
