package database

import (
	"goapi-starter/internal/logger"
	"goapi-starter/internal/metrics"
	"time"

	"gorm.io/gorm"
)

// AddMetricsCallbacks adds callbacks to GORM to track database operations
func AddMetricsCallbacks() {
	logger.Debug().Msg("Adding database metrics callbacks")

	// Create operation
	DB.Callback().Create().Before("gorm:create").Register("metrics:create_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
		logger.Debug().
			Str("entity", getEntityName(db)).
			Msg("Starting database create operation")
	})

	DB.Callback().Create().After("gorm:create").Register("metrics:create_after", func(db *gorm.DB) {
		entity := getEntityName(db)
		metrics.RecordDatabaseOperation("create", entity)
		logger.Debug().
			Str("entity", entity).
			Bool("error", db.Error != nil).
			Msg("Completed database create operation")
	})

	// Query operation
	DB.Callback().Query().Before("gorm:query").Register("metrics:query_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
		logger.Debug().
			Str("entity", getEntityName(db)).
			Str("query", db.Statement.SQL.String()).
			Msg("Starting database query operation")
	})

	DB.Callback().Query().After("gorm:query").Register("metrics:query_after", func(db *gorm.DB) {
		entity := getEntityName(db)
		metrics.RecordDatabaseOperation("query", entity)
		logger.Debug().
			Str("entity", entity).
			Bool("error", db.Error != nil).
			Int64("rows_affected", db.RowsAffected).
			Msg("Completed database query operation")
	})

	// Update operation
	DB.Callback().Update().Before("gorm:update").Register("metrics:update_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
		logger.Debug().
			Str("entity", getEntityName(db)).
			Msg("Starting database update operation")
	})

	DB.Callback().Update().After("gorm:update").Register("metrics:update_after", func(db *gorm.DB) {
		entity := getEntityName(db)
		metrics.RecordDatabaseOperation("update", entity)
		logger.Debug().
			Str("entity", entity).
			Bool("error", db.Error != nil).
			Int64("rows_affected", db.RowsAffected).
			Msg("Completed database update operation")
	})

	// Delete operation
	DB.Callback().Delete().Before("gorm:delete").Register("metrics:delete_before", func(db *gorm.DB) {
		db.Set("metrics:start_time", time.Now())
		logger.Debug().
			Str("entity", getEntityName(db)).
			Msg("Starting database delete operation")
	})

	DB.Callback().Delete().After("gorm:delete").Register("metrics:delete_after", func(db *gorm.DB) {
		entity := getEntityName(db)
		metrics.RecordDatabaseOperation("delete", entity)
		logger.Debug().
			Str("entity", entity).
			Bool("error", db.Error != nil).
			Int64("rows_affected", db.RowsAffected).
			Msg("Completed database delete operation")
	})

	logger.Info().Msg("Database metrics callbacks registered successfully")
}

// getEntityName tries to determine the entity name from the GORM DB context
func getEntityName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	return "unknown"
}
