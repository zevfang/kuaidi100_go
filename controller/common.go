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
	"github.com/ahmetb/go-linq"
	"kuaidi100_go/log"
	"strings"
)

/***************************************推送处理******************************************/
/*
	接受推送

	200: 提交成功
	500: 服务器错误
	其他错误请自行定义
	100 参数错误
	300 签名错误
	400 内容解析错误或非法数据

*/

//func SetLjTime() {
//	order, err := model.GetTop()
//	if err != nil {
//		fmt.Println(err)
//	}
//	for i := 0; i < len(order); i++ {
//
//		d := order[i].Data
//		var callBackData model.CallPostData
//		err := json.Unmarshal([]byte(d), &callBackData)
//		if err != nil {
//			fmt.Println("转换错误")
//			continue
//		}
//
//		order[i].LjTime = getStatusMinTime(callBackData.LastResult.Data, "揽件")
//		order[i].ZtTime = getStatusMinTime(callBackData.LastResult.Data, "在途")
//
//		model.UpdTop(order[i])
//		fmt.Printf("nu:%s , index:%d \r\n ", order[i].Nu, i)
//	}
//	if len(order) == 0 {
//		fmt.Println("完成，请查询检查一下")
//	}
//}

func CallBack(c *gin.Context) {
	param := c.PostForm("param")
	sign := c.PostForm("sign")
	//验证
	if param == "" || sign == "" {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "100",
			"message":    "参数错误",
		})
		return
	}

	//验签
	signStr := system.Md5(param + system.GetConfiguration().Salt)
	if strings.ToUpper(sign) != strings.ToUpper(signStr) {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "300",
			"message":    "签名验证错误",
		})
		return
	}

	//解析数据
	var callBackData model.CallPostData
	err := json.Unmarshal([]byte(param), &callBackData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "400",
			"message":    "内容解析错误或非法数据",
		})
		return
	}

	kdOrder := model.KdOrder{
		KdStatus:     callBackData.Status,
		KdStatusName: converToKdStatusName(callBackData.Status),
		KdMessage:    callBackData.Message,
		State:        callBackData.LastResult.State,
		StateName:    convertToStateName(callBackData.LastResult.State),
		Com:          callBackData.LastResult.Com,
		Shelp:        "",
		Nu:           callBackData.LastResult.Nu,
		Data:         param,
		ZtName:       "在途",
		ZtTime:       getStatusMinTime(callBackData.LastResult.Data, "在途"),
		LjName:       "揽件",
		LjTime:       getStatusMinTime(callBackData.LastResult.Data, "揽件"),
		YnName:       "疑难",
		YnTime:       getStatusMinTime(callBackData.LastResult.Data, "疑难"),
		QsName:       "签收",
		QsTime:       getStatusMinTime(callBackData.LastResult.Data, "签收"),
		TqName:       "退签",
		TqTime:       getStatusMinTime(callBackData.LastResult.Data, "退签"),
		PjName:       "派件",
		PjTime:       getStatusMinTime(callBackData.LastResult.Data, "派件"),
		ThName:       "退回",
		ThTime:       getStatusMinTime(callBackData.LastResult.Data, "退回"),
		Created:      "",
		Updated:      "",
	}

	//查看是否存在(orderErr：如果结果集是空的，返回一个错误)
	order, orderErr := model.GetKdOrderOne(callBackData.LastResult.Nu)

	var errs error
	//新增数据
	if order.Nu == "" && orderErr != nil {
		//获取阿芙快递商编码
		shelp, err := model.GetTradeOrderShelp(callBackData.LastResult.Nu)
		if err != nil {
			log.Log.Info(fmt.Sprintf("快递商编码查询异常:[%s]%s", callBackData.LastResult.Nu, err))
		}
		//新增
		kdOrder.Shelp = shelp
		kdOrder.Created = system.GetNow()
		kdOrder.Updated = system.GetNow()
		errs = model.InsertKdOrder(kdOrder)
	} else {
		//修改
		kdOrder.Shelp = order.Shelp
		kdOrder.Created = order.Created
		kdOrder.Updated = system.GetNow()
		errs = model.UpdateKdOrder(kdOrder)
	}
	if errs != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":     false,
			"returnCode": "500",
			"message":    "服务器错误,数据未被存储" + errs.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result":     true,
		"returnCode": "200",
		"message":    "提交成功",
	})
}

// 获取流转节点时间
func getStatusMinTime(r []model.CallPostLastResultData, status string) string {
	data_cout := len(r)
	if data_cout > 0 {
		var minTime string = ""
		var isfalg bool
		lastData := linq.From(r).WhereT(func(c model.CallPostLastResultData) bool {
			return c.Status == status
		}).Results()

		//查询关键字是否存在
		if len(lastData) > 0 {
			isfalg = true
		}

		//如果存在关键字
		if isfalg {
			//默认列表首条数据
			t := linq.From(lastData).SelectT(func(c model.CallPostLastResultData) string {
				return c.Ftime
			}).Min()
			minTime = fmt.Sprintf("%s", t)

			//如果在途数量大于1，则查询在途中的第二条，否则选择

			if status == "在途" && len(lastData) > 1 {
				z := lastData[len(lastData)-2]
				minTime = z.(model.CallPostLastResultData).Ftime
			}

		} else {
			//从在途里面取首条
			if status == "揽件" {
				d := linq.From(r).WhereT(func(c model.CallPostLastResultData) bool {
					return c.Status == "在途"
				}).SelectT(func(c model.CallPostLastResultData) string {
					return c.Ftime
				}).Min()
				minTime = fmt.Sprintf("%s", d)
			}

			if status == "在途" && len(lastData) == 0 {
				minTime = r[len(r)-1].Ftime
			}
		}
		return minTime
	}
	return ""
}

func converToKdStatusName(kdStatus string) string {
	var kdStatusName string
	switch kdStatus {
	case "polling":
		kdStatusName = "监控中"
		break
	case "shutdown":
		kdStatusName = "结束"
		break
	case "abort":
		kdStatusName = "中止"
		break
	case "updateall":
		kdStatusName = "重新推送"
		break
	default:
		kdStatusName = "传入状态错误"
		break
	}
	return kdStatusName
}

func convertToStateName(state string) string {
	var statusName string
	switch state {
	case "0":
		statusName = "在途"
		break
	case "1":
		statusName = "揽件"
		break
	case "2":
		statusName = "疑难"
		break
	case "3":
		statusName = "签收"
		break
	case "4":
		statusName = "退签"
		break
	case "5":
		statusName = "派件"
		break
	case "6":
		statusName = "退回"
		break
	default:
		statusName = "传入状态错误"
		break
	}
	return statusName
}

/***************************************************订阅处理*********************************************/

/*
	订阅数据

	200: 提交成功
	701: 拒绝订阅的快递公司
	700: 订阅方的订阅数据存在错误（如不支持的快递公司、单号为空、单号超长等）或错误的回调地址
	702: POLL:识别不到该单号对应的快递公司
	600: 您不是合法的订阅者（即授权Key出错）
	601: POLL:KEY已过期
	500: 服务器错误（即快递100的服务器出理间隙或临时性异常，有时如果因为不按规范提交请求，比如快递公司参数写错等，也会报此错误）
	501: 重复订阅（请格外注意，501表示这张单已经订阅成功且目前还在跟踪过程中（即单号的status=polling），快递100的服务器会因此忽略您最新的此次订阅请求，从而返回501。一个运单号只要提交一次订阅即可，若要提交多次订阅，请在收到单号的status=abort或shutdown后隔半小时再提交订阅
*/

func PollOrder() {
	// 获取待订阅数据
	data, err := model.GetTopOrder()
	if err != nil {
		fmt.Println(err)
	}
	if len(data) == 0 {
		return
	}
	pollData := converToCom(data)
	fmt.Printf("订阅开始:%d", len(pollData))
	// 循环订阅数据
	for _, v := range pollData {

		resultData, err := postData(system.GetConfiguration().SubscribeUrl, v)
		if err != nil {
			log.Log.Error(fmt.Sprintf("订阅发生错误：%s", err))
		}

		sLog := model.KdSubscribeLog{
			LogisticsOrder: v.Param.Number,
			Result:         resultData.Result,
			ReturnCode:     resultData.ReturnCode,
			Message:        resultData.Message,
		}

		//返回值处理
		err = model.TranOrderAndResult(sLog)
		if err != nil {
			log.Log.Error(fmt.Sprintf("对账数据插入失败:%s", err))
		}
	}
}

// 并发访问（暂时不启用）
func asyncPollOrder(c *gin.Context) {

	// 获取待订阅数据

	data, err := model.GetTopOrder()
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
	//	model.TranOrderAndResult()
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
