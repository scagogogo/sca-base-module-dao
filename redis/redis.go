package redis

import (
	"github.com/go-redis/redis"
	"github.com/golang-infrastructure/go-pointer"
	sca_base_module_config "github.com/scagogogo/sca-base-module-config"
)

var Client *redis.Client

func init() {

	if sca_base_module_config.Config == nil {
		return
	}

	if !pointer.FromPointerOrDefault(sca_base_module_config.Config.Redis.AutoInit, false) {
		return
	}

	InitRedis(sca_base_module_config.Config)
}

func InitRedis(config *sca_base_module_config.Configuration) {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Passwd,
	})
}
