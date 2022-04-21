package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

var ViperConf ViperConfig
var fileName = "../dev-config.yaml"

func init() {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.ReadInConfig()
	err := v.Unmarshal(&ViperConf)
	if err != nil {
		fmt.Println()
	}
	fmt.Println(ViperConf)
	fmt.Println("Viper初始化完成。。。")
	InitRedis()
}

type ViperConfig struct {
	RedisConfig      RedisConfig      `mapstructure:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul"`
	AccountSrvConfig AccountSrvConfig `mapstructure:"account_srv"`
	AccountWebConfig AccountWebConfig `mapstructure:"account_web"`
}
