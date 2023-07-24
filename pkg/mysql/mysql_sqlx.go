package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
)

var SqlX *sqlx.DB

func InitSqlX(config *sca_base_module_config.Configuration) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if config.Database.Mysql.DSN == "" {
		return fmt.Errorf("mysql dns is empty")
	}
	db, err := sqlx.Connect("mysql", config.Database.Mysql.DSN)
	if err != nil {
		return err
	}
	SqlX = db
	return nil
}