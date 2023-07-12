package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-infrastructure/go-pointer"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
)

func init() {

	// 配置文件都没有，那就没得玩了
	if sca_base_module_config.Config == nil {
		return
	}

	// 除非手动开启了自动初始化mysql模块，否则不自动初始化
	if !pointer.FromPointerOrDefault(sca_base_module_config.Config.Database.Mysql.AutoInit, false) {
		return
	}

	// 初始化SqlX，如果需要的话
	if pointer.FromPointerOrDefault(sca_base_module_config.Config.Database.Mysql.Driver.SqlX, false) {
		err := InitSqlX(sca_base_module_config.Config)
		if err != nil {
			panic(err)
		}
	}

	// 初始化Gorm，如果需要的话
	if pointer.FromPointerOrDefault(sca_base_module_config.Config.Database.Mysql.Driver.Gorm, false) {
		err := InitGorm(sca_base_module_config.Config)
		if err != nil {
			panic(err)
		}
	}

}
