package database

import (
	"goapi-starter/internal/metrics"
	"time"

	"gorm.io/gorm"
)

// AddMetricsCallbacks adds callbacks to GORM to track database operations
func AddMetricsCallbacks() {
	// Create operation
	DB.Callback().Create().Before("gorm:create").Register("metrics:create_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
	})

	DB.Callback().Create().After("gorm:create").Register("metrics:create_after", func(db *gorm.DB) {
		metrics.RecordDatabaseOperation("create", getEntityName(db))
	})

	// Query operation
	DB.Callback().Query().Before("gorm:query").Register("metrics:query_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
	})

	DB.Callback().Query().After("gorm:query").Register("metrics:query_after", func(db *gorm.DB) {
		metrics.RecordDatabaseOperation("query", getEntityName(db))
	})

	// Update operation
	DB.Callback().Update().Before("gorm:update").Register("metrics:update_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
	})

	DB.Callback().Update().After("gorm:update").Register("metrics:update_after", func(db *gorm.DB) {
		metrics.RecordDatabaseOperation("update", getEntityName(db))
	})

	// Delete operation
	DB.Callback().Delete().Before("gorm:delete").Register("metrics:delete_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
	})

	DB.Callback().Delete().After("gorm:delete").Register("metrics:delete_after", func(db *gorm.DB) {
		metrics.RecordDatabaseOperation("delete", getEntityName(db))
	})
}

// getEntityName tries to determine the entity name from the GORM DB context
func getEntityName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	return "unknown"
}
