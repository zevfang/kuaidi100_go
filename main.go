package main

import (
	"os"
	"fmt"
	"github.com/kardianos/service"
	"github.com/gin-gonic/gin"
	"kuaidi100_go/log"
	"kuaidi100_go/controller"
	"net/http"
	"kuaidi100_go/system"
	"github.com/claudiu/gocron"
	"kuaidi100_go/model"
)

type daemon struct{}

func (p *daemon) Start(s service.Service) error {
	go p.Run()
	return nil
}

func (p *daemon) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *daemon) Run() {

	//初始化log
	log.NewLogger()
	log.Log.Info("app init \n\r")

	configFilePath := fmt.Sprintf("%s%s", system.GetCurrentDirectory(), "/conf/config.ini")
	convertFilePath := fmt.Sprintf("%s%s", system.GetCurrentDirectory(), "/conf/convert.json")

	// 加载配置文件
	if err := system.LoadConfiguration(configFilePath); err != nil {
		log.Log.Error(err.Error())
		fmt.Println(err)
		return
	}

	// 加载快递编码对照表
	if err := system.LoadComs(convertFilePath); err != nil {
		log.Log.Error(err.Error())
		fmt.Println(err)
		return
	}

	// 注册DB
	if err := model.InitModel(); err != nil {
		log.Log.Error(err.Error())
		fmt.Println(err)
		return
	}

	gin.SetMode(system.GetConfiguration().Env)
	router := gin.Default()
	router.Use(gin.ErrorLogger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "快递100服务",
		})
	})

	//接受推送
	router.POST("/callback", controller.CallBack)
	fmt.Println("接受推送服务已开启")
	//定时订阅
	if system.GetConfiguration().PollState == 1 {
		//获取订阅间隔
		m := uint64(system.GetConfiguration().PollMinutes)

		gocron.Every(m).Seconds().Do(controller.PollOrder)
		gocron.Start()
		fmt.Println("注册订阅服务已开启")
	}

	router.Run(system.GetConfiguration().Addr)
}

func main() {

	svcConfig := &service.Config{
		Name:        "快递信息订阅推送",
		DisplayName: "kuidi100_go",
		Description: "快递100获取物流信息的订阅/推送服务",
	}

	prg := &daemon{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println("Create service error => ", err)
	}
	if len(os.Args) > 1 {

		if os.Args[1] == "install" {
			err := s.Install()
			if err != nil {
				fmt.Println("Install service error=> ", err)
				os.Exit(1)
			}
			fmt.Println("Successful install services")
			return
		}

		if os.Args[1] == "remove" {
			err := s.Uninstall()
			if err != nil {
				fmt.Println("Remove service error=> ", err)
				os.Exit(1)
			}
			fmt.Println("Successful remove services")
			return
		}

		if os.Args[1] == "restart" {
			err := s.Restart()
			if err != nil {
				fmt.Println("Restart service error=> ", err)
				os.Exit(1)
			}
			fmt.Println("Successful restart services")
			return
		}
	}

	err = s.Run()
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}

}
