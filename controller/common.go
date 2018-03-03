package controller

import (
	"kuaidi100_go/model"
	"kuaidi100_go/system"
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"sync"
)

/*
	接受推送

	200: 提交成功
	500: 服务器错误
	其他错误请自行定义
	100 参数错误
	300 签名错误
	400 内容解析错误或非法数据

*/
func CallBack(c *gin.Context) {
	param := c.PostForm("param")
	sign := c.PostForm("sign")
	//验证
	if param != "" && sign != "" {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "100",
			"message":    "参数错误",
		})
		return
	}
	//验签
	signStr := system.Md5(param + system.GetConfiguration().Salt)
	if sign != signStr {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "300",
			"message":    "签名验证错误",
		})
		return
	}

	//解析数据
	var call model.CallPostData
	err := json.Unmarshal([]byte(param), &call)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "400",
			"message":    "内容解析错误或非法数据",
		})
		return
	}

	//组装数据准备入库

}

func ConverToStatus(callData model.CallPostData) error {
	return nil
}

/*
	订阅数据

	200: 提交成功
	701: 拒绝订阅的快递公司
	700: 订阅方的订阅数据存在错误（如不支持的快递公司、单号为空、单号超长等）或错误的回调地址
	702: POLL:识别不到该单号对应的快递公司
	600: 您不是合法的订阅者（即授权Key出错）
	601: POLL:KEY已过期
	500: 服务器错误（即快递100的服务器出理间隙或临时性异常，有时如果因为不按规范提交请求，比如快递公司参数写错等，也会报此错误）
	501:重复订阅（请格外注意，501表示这张单已经订阅成功且目前还在跟踪过程中（即单号的status=polling），快递100的服务器会因此忽略您最新的此次订阅请求，从而返回501。一个运单号只要提交一次订阅即可，若要提交多次订阅，请在收到单号的status=abort或shutdown后隔半小时再提交订阅
*/

func PollOrder(c *gin.Context) {

	// 获取待订阅数据
	t := new(model.TradeOrder)
	data, err := t.GetTopOrder()
	if err != nil {
		fmt.Println(err)
	}
	pollData := converToCom(data)
	//var results []model.ResultData
	//for _,v:=range  pollData{
	//
	//	result, err := PostData(system.GetConfiguration().SubscribeUrl, v)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	results = append(results,result)
	//}
	//for _, v := range results {
	//	t.TranOrderAndResult()
	//}

	//返回
	c.JSON(http.StatusOK, pollData)
}

// 并发访问（暂时不启用）
func AsyncPollOrder(c *gin.Context) {

	// 获取待订阅数据
	t := new(model.TradeOrder)
	data, err := t.GetTopOrder()
	if err != nil {
		fmt.Println(err)
	}
	pollData := converToCom(data)

	//带锁的Map
	var resultLogs = struct {
		sync.RWMutex
		Data map[string]model.ResultData
	}{Data: make(map[string]model.ResultData)}

	var wg sync.WaitGroup
	for i := 0; i < len(pollData); i++ {
		wg.Add(1)
		go func(v model.PollPostData) {
			defer wg.Done()
			//result, err := postData(system.GetConfiguration().SubscribeUrl, v)
			//if err != nil {
			//	fmt.Println(err)
			//}
			result := model.ResultData{
				true,
				"200",
				"成功" + v.Param.Number,
			}
			//存储返回值
			resultLogs.Lock()
			defer resultLogs.Unlock()
			resultLogs.Data[v.Param.Number] = result
		}(pollData[i])
	}
	wg.Wait()

	//for _, v := range results {
	//	t.TranOrderAndResult()
	//}

	//返回
	c.JSON(http.StatusOK, gin.H{
		"pcount": len(pollData),
		"rcount": len(resultLogs.Data),
		"data":   resultLogs,
	})
}

//请求
func postData(requestUrl string, data model.PollPostData) (model.ResultData, error) {
	var result model.ResultData
	//转换字符
	b, err := json.Marshal(data.Param)
	if err != nil {
		return result, errors.New("convert json error.")
	}
	paramStr := string(b)

	//发起请求
	res, err := http.PostForm(requestUrl, url.Values{
		"schema": {data.Schema},
		"param":  {paramStr},
	})
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	//获取返回
	err = json.NewDecoder(res.Body).Decode(&result)
	return result, err
}

//转换并匹配快递100参数code
func converToCom(tradeOrders []model.TradeOrder) []model.PollPostData {
	var poll_arr []model.PollPostData
	for _, v := range tradeOrders {
		var m model.PollPostData
		m.Schema = "json"

		companyCode := system.GetComByCodeName(v.Shelp, v.LogisticsComp)

		m.Param.Company = companyCode
		m.Param.Number = v.LogisticsOrder
		m.Param.Key = system.GetConfiguration().Key

		m.Param.Parameters.CallbackUrl = system.GetConfiguration().CallbackUrl
		m.Param.Parameters.Salt = system.GetConfiguration().Salt
		//添加此字段表示开通行政区域解析功能 返回（status、areaCode、areaName）
		m.Param.Parameters.Resultv2 = "1"
		//如果没有匹配到快递商，则加入智能识别
		if companyCode == "" {
			m.Param.Parameters.AutoCom = "1"
		}
		poll_arr = append(poll_arr, m)
	}
	return poll_arr
}
