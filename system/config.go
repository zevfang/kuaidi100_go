package system

import (
	"github.com/go-ini/ini"
)

type Configuration struct {
	Addr           string `ini:"addr"`
	Key            string `ini:"app_key"`          //订阅推送接口Key（由快递100颁发）
	SubscribeUrl   string `ini:"subscribe_url"`    //订阅地址
	CallbackUrl    string `ini:"callback_url"`     //回调地址
	RedisMaxIdle   int    `ini:"redis_max_idle"`   //最大的空闲连接数
	RedisMaxActive int    `ini:"redis_max_active"` //最大的激活连接数(同时最多有N个连接)
	RedisHost      string `ini:"redis_host"`
	RedisDb        string `ini:"redis_db"`
	RedisPassWord  string `ini:"redis_pass_word"`
	MsSqlServer    string `ini:"mssql_server"`
	MsSqlDataBase  string `ini:"mssql_db"`
	MsSqlPort      int    `ini:"mssql_port"`
	MsSqlUid       string `ini:"mssql_uid"`
	MsSqlPwd       string `ini:"mssql_pwd"`
}

var configuration Configuration

func LoadConfiguration(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	return cfg.MapTo(&configuration)
}

func GetConfiguration() *Configuration {
	return &configuration
}
