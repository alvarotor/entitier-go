package mocks

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupGORMSqlite(t *testing.T, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			t.Fatalf("failed to migrate database schema for model %v: %v", model, err)
		}
	}

	return db
}
