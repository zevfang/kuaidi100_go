package model

import (
	"github.com/garyburd/redigo/redis"
	"kuaidi100_go/system"
	"time"
)

type Redis struct {
}


var (
	// 定义常量
	RedisClient *redis.Pool
)

func InitRedis() {
	//参数变量
	MAX_IDLE := system.GetConfiguration().RedisMaxIdle
	MAX_ACTIVE := system.GetConfiguration().RedisMaxActive
	HOST := system.GetConfiguration().RedisHost
	DB := system.GetConfiguration().RedisDb
	PASS_WORD := system.GetConfiguration().RedisPassWord
	//建立连接池
	RedisClient = &redis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: 180 * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", HOST)
			if err != nil {
				return nil, err
			}
			//密码
			if _, err := c.Do("AUTH", PASS_WORD); err != nil {
				c.Close()
				return nil, err
			}
			// 选择db
			if _, err := c.Do("SELECT", DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

//func (r *Redis) SetExpress(s string) error {
//	conn,err:=RedisClient.Dial()
//	if err!=nil {
//		return  err
//	}
//	conn.Flush().Error()
//}