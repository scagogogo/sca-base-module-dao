package mysql

import (
	"fmt"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Gorm *gorm.DB

func InitGorm(config *sca_base_module_config.Configuration) error {

	if config == nil {
		return fmt.Errorf("config nil")
	}

	if config.Database.Mysql.DSN == "" {
		return fmt.Errorf("dsn is empty")
	}

	db, err := NewMySQL(config.Database.Mysql.DSN)
	if err != nil {
		return err
	}

	// connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db: %w", err)
	}
	// need same open connections and idle connections
	// @see: https://github.com/go-sql-driver/mysql/issues/991#issuecomment-526035935
	if config.Database.Mysql.Connection.MaxIdle != nil {
		sqlDB.SetMaxIdleConns(*config.Database.Mysql.Connection.MaxIdle)
	}

	if config.Database.Mysql.Connection.MaxOpen != nil {
		sqlDB.SetMaxOpenConns(*config.Database.Mysql.Connection.MaxOpen)
	}

	if config.Database.Mysql.Connection.MaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(*config.Database.Mysql.Connection.MaxLifetime)
	}

	Gorm = db
	return nil
}

func NewMySQL(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
}

func TableWithSuffix(t schema.Tabler, suffix string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(t.TableName() + "_" + suffix)
	}
}
