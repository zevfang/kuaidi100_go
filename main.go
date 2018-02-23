package main

import (
	"fmt"
	"github.com/claudiu/gocron"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"kuaidi100_go/model"
	"kuaidi100_go/system"
	"net/http"
)

func main() {

	// 加载配置文件
	if err := system.LoadConfiguration("conf/config.ini"); err != nil {
		fmt.Println(err)
		return
	}

	// 加载快递编码对照表
	if err := system.LoadComs("conf/convert.json"); err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()
	router.Use(gin.ErrorLogger())
	router.GET("/getcom", func(c *gin.Context) {
		c.JSON(200, system.GetComArray())
	})
	router.GET("/getdata", func(c *gin.Context) {
		t := new(model.TradeOrder)
		res, err := t.GetTopOrder()
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK,res)
	})
	gocron.Every(1).Day().Do(func() {
		fmt.Println("hello")
	})
	gocron.Every(1).Hour().Do(func() {
		fmt.Println("nihao")
	})
	gocron.Start()

	router.Run(system.GetConfiguration().Addr)

}

func GetKD100Json(url string) {
	res, err := http.Get("http://www.baidu.com")
	if err != nil {
		//log.Error("http err is get baidu")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println(string(b))
}
