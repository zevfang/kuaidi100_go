package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/garyburd/redigo/redis"
	"github.com/go-ini/ini"
)

type Com struct {
	ACode string
	KCode string
	KName string
}

type Config struct {
	/*	是
		身份授权key，请 快递查询接口 进行申请（大小写敏感）
	*/
	Id string `ini:"app_key"`
	/*	是
		要查询的快递公司代码，不支持中文，对应的公司代码见《API URL 所支持的快递公司及参数说明》和《支持的国际类快递及参数说明》。
		如果找不到您所需的公司，请发邮件至 kuaidi@kingdee.com 咨询（大小写不敏感）
	*/
	Com string `ini:"com"`
	/*	是
		要查询的快递单号，请勿带特殊符号，不支持中文（大小写不敏感）
	*/
	Nu string `ini:"nu"`
	/*	是
		已弃用字段，无意义，请忽略。
	*/
	Valicode string `ini:"valicode"`
	/*	是
		返回类型：
			0：返回json字符串，
			1：返回xml对象，
			2：返回html对象，
			3：返回text文本。
		如果不填，默认返回json字符串。
	*/
	Show string `ini:"show"`
	/*	是
		返回信息数量：
			1:返回多行完整的信息，
			0:只返回一行信息。
		不填默认返回多行。
	*/
	Muti string `ini:"muti"`
	/*	是
		排序：
			desc：按时间由新到旧排列，
			asc ：按时间由旧到新排列。
		不填默认返回倒序（大小写不敏感）
	*/
	Order string `ini:"order"`
}

var log *logger.Logger
var cfg = new(Config)
var coms *[]Com

var cfg_redis *ini.Section
var redisClient *redis.Pool

func InitLogger() {
	var (
		err     error
		logFile *os.File
	)
	logFile, err = os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file error.")
	}
	log, err = logger.New("main", 1, logFile)
	if err != nil {
		panic("log init error.")
	}
}

func InitKuaiDi() {
	file, err := os.OpenFile("kuaidi.json", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("kuaidi.json err")
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("kuaidi.json err")
	}

	err = json.Unmarshal([]byte(b), &coms)
	if err != nil {
		fmt.Println(err)
	}
}

func InitKuaiDi100Config() {
	c, err := ini.Load("config.ini")
	if err != nil {
		log.Error("config is error")
	}
	err = c.Section("kuaidi100").MapTo(cfg)
	if err != nil {
		log.Error("config is node kuaidi100 section error")
	}
}

func InitRedisConfig() {
	c, err := ini.Load("config.ini")
	if err != nil {
		log.Error("config is error")
	}
	cfg_redis = c.Section("redis")
}

func InitRedis() {
	MAX_IDLE, err := cfg_redis.Key("redis_maxidle").Int()
	if err != nil {
		log.Error("Type conversion error")
	}
	MAX_ACTIVE, err := cfg_redis.Key("redis_maxactive").Int()
	if err != nil {
		log.Error("Type conversion error")
	}

	HOST := cfg_redis.Key("redis_host").String()
	DB := cfg_redis.Key("redis_db").String()
	PASS_WORD := cfg_redis.Key("redis_pass_word").String()
	//建立连接池
	redisClient = &redis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", HOST)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			//密码
			if _, err := c.Do("AUTH", PASS_WORD); err != nil {
				fmt.Println(err)
				c.Close()
				return nil, err
			}
			// 选择db
			if _, err := c.Do("SELECT", DB); err != nil {
				fmt.Println(err)
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func init() {
	InitLogger()
	InitKuaiDi()
	InitKuaiDi100Config()
	InitRedisConfig()
	InitRedis()
}

func main() {

	// fmt.Println(cfg)
	// fmt.Println(coms)
	// log.Error("这是个错误")
	//fmt.Println(cfg_redis)
	c, err := redisClient.Dial()
	if err != nil {
		fmt.Println(err)
	}

	_, err = c.Do("SET", "express:com", "韵达")
	
	if err != nil {
		log.Error("这是个错误")
	}
	defer c.Close()
	//apiurl := fmt.Sprintf("http://api.kuaidi100.com/api?id=%s&com=[]&nu=[]&valicode=[]&show=[0|1|2|3]&muti=[0|1]&order=[desc|asc]"
	// GetKD100Json(url)
}

func GetKD100Json(url string) {
	res, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Error("http err is get baidu")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println(string(b))
}
