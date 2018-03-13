package system

import (
	"github.com/go-ini/ini"
)

type Configuration struct {
	Env                string `ini:"env"`
	Addr               string `ini:"addr"`
	Key                string `ini:"app_key"`       //订阅推送接口Key（由快递100颁发）
	SubscribeUrl       string `ini:"subscribe_url"` //订阅地址
	CallbackUrl        string `ini:"callback_url"`  //回调地址
	Salt               string `ini:"salt"`
	RedisMaxIdle       int    `ini:"redis_max_idle"`   //最大的空闲连接数
	RedisMaxActive     int    `ini:"redis_max_active"` //最大的激活连接数(同时最多有N个连接)
	RedisHost          string `ini:"redis_host"`
	RedisDb            int `ini:"redis_db"`
	RedisPassWord      string `ini:"redis_pass_word"`
	MsSqlServer        string `ini:"mssql_server"`
	MsSqlDataBase      string `ini:"mssql_db"`
	MsSqlPort          int    `ini:"mssql_port"`
	MsSqlUid           string `ini:"mssql_uid"`
	MsSqlPwd           string `ini:"mssql_pwd"`
	ServiceName        string `ini:"service_name"`         //显示名称
	ServiceDisplayName string `ini:"service_display_name"` //服务名称
	ServiceDescription string `ini:"service_description"`  //备注信息
	PollState          int    `ini:"poll_state"`           //是否开启订阅
	PollMinutes        int    `ini:"poll_minutes"`         //订阅间隔（秒）
	OncePollCount      int    `ini:"once_poll_count"`      //每次订阅的数量
	StartDate          string `ini:"start_date"`           //开始查询时间
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
