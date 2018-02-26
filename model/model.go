package model

func InitModel() error {
	var err error
	//初始化Redis
	InitRedis()

	//初始SqlServer
	err = InitMsSql()
	return err
}

// 源数据结构体
type TradeOrder struct {
	Shelp          string `db:"shelp"`
	LogisticsComp  string `db:"logistics_comp"`
	LogisticsOrder string `db:"logistics_order"`
}

/*
 	schema= json
    param={
        "company":"ems",
        "number":"em263999513jp",
        "from":"广东省深圳市南山区",
        "to":"北京市朝阳区",
        "key":"XXX ",
        "parameters":{
            "callbackurl":"您的回调接口的地址，如http://www.您的域名.com/kuaidi?callbackid=...",
            "salt":"XXXXXXXXXX",
            "resultv2":"1",
            "autoCom":"1",
            "interCom"："1",
            "departureCountry":"CN",
            "departureCom":"ems",
            "destinationCountry":"JP",
            "destinationCom":"japanposten"
        }
    }
*/
// 订阅参数结构
type PollPostData struct {
	Schema string `json:"schema"`
	Param struct {
		Company string `json:"company"`
		Number  string `json:"number"`
		Key     string `json:"key"`
		Parameters struct {
			CallbackUrl string `json:"callbackurl"`
			Salt        string `json:"salt"` //回调使用（签名字符串，随机即可）
			AutoCom     string `json:"autoCom"`
		} `json:"parameters"`
	} `json:"param"`
}

/*
	{
		"status": "polling",
		"billstatus": "got",
		"message": "",
		"autoCheck": "1",
		"comOld": "yuantong",
		"comNew": "ems",
		"lastResult": {
			"message": "ok",
			"state": "0",
			"status": "200",
			"condition": "F00",
			"ischeck": "0",
			"com": "yuantong",
			"nu": "V030344422",
			"data": [{
				"context": "上海分拨中心/装件入车扫描 ",
				"time": "2012-08-28 16:33:19",
				"ftime": "2012-08-28 16:33:19",
				"status": "在途",
				"areaCode": "310000000000",
				"areaName": "上海市"
			}, {
				"context": "上海分拨中心/下车扫描 ",
				"time": "2012-08-27 23:22:42",
				"ftime": "2012-08-27 23:22:42",
				"status": "在途",
				"areaCode": "310000000000",
				"areaName": "上海市"

			}]
		}
	}
*/
//推送数据结构
type CallPostData struct {
}

/*
	{
        "result":true,
        "returnCode":"200",
        "message":"提交成功"
    }
*/
//返回结构
type ResultData struct {
	Result     bool
	ReturnCode int
	Message    string
}
