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


// 订阅参数结构
type PollPostData struct {
	Schema string `json:"schema"`
	Param struct {
		Company string `json:"company"`
		Number  string `json:"number"`
		Key     string `json:"key"`
		Parameters struct {
			CallbackUrl string `json:"callbackurl"`
			Salt        string `json:"salt"` //回调使用（签名字符串，随机即可---配置内base64加密码）
			AutoCom     string `json:"autoCom"`
			Resultv2    string `json:"resultv2"` //值为1 填写之后推送返回（status、areaCode、areaName）
		} `json:"parameters"`
	} `json:"param"`
}

//推送接受数据结构
type CallPostData struct {
	Status     string `json:"status"`
	BillStatus string `json:"billstatus"`
	Message    string `json:"message"`
	AutoCheck  string `json:"autoCheck"`
	ComOld     string `json:"comOld"`
	ComNew     string `json:"comNew"`
	LastResult struct {
		Message   string                   `json:"message"`
		State     string                   `json:"state"`
		Status    string                   `json:"status"`
		Condition string                   `json:"condition"`
		Ischeck   string                   `json:"ischeck"`
		Com       string                   `json:"com"`
		Nu        string                   `json:"nu"`
		Data      []CallPostLastResultData `json:"data"`
	} `json:"lastResult"`
}

type CallPostLastResultData struct {
	Context  string `json:"context"`
	Time     string `json:"time"`
	Ftime    string `json:"ftime"`
	Status   string `json:"status"`
	AreaCode string `json:"areaCode"`
	AreaName string `json:"areaName"`
}

// 推送写入库结构

type KdOrder struct {
	Id           int64  `json:"id" db:"id"`                         //自增编号
	KdStatus     string `json:"kd_status" db:"kd_status"`           //监控状态:polling:监控中，shutdown:结束，abort:中止，updateall：重新推送
	KdStatusName string `json:"kd_status_name" db:"kd_status_name"` //转换后的名称
	KdMessage    string `json:"kd_message" db:"kd_message"`         //监控状态相关消息，如:3天查询无记录，60天无变化  =>abort:中止
	State        string `json:"state" db:"state"`                   //当前签收状态，包括0在途中、1已揽收、2疑难、3已签收、4退签、5同城派送中、6退回等状态
	StateName    string `json:"state_name" db:"state_name"`         //当前快递状态名称
	Com          string `json:"com" db:"com"`                       //（kuaidi100）快递公司编码
	Shelp        string `json:"shelp" db:"shelp"`                   //（本地）快递公司编码 （需要转换后入库）
	Nu           string `json:"nu" db:"nu"`                         //快递单号
	Data         string `json:"data" db:"data"`                     //全量快递信息 json字符串
	ZtName       string `json:"zt_name" db:"zt_name"`               //0	在途	快件处于运输过程中
	ZtTime       string `json:"zt_time" db:"zt_time"`
	LjName       string `json:"lj_name" db:"lj_name"` 				//1	揽件	快件已由快递公司揽收
	LjTime       string `json:"lj_time" db:"lj_time"`
	YnName       string `json:"yn_name" db:"yn_name"` 				//2	疑难	快递100无法解析的状态，或者是需要人工介入的状态，比方说收件人电话错误。
	YnTime       string `json:"yn_time" db:"yn_time"`
	QsName       string `json:"qs_name" db:"qs_name"` 				//3	签收	正常签收
	QsTime       string `json:"qs_time" db:"qs_time"`
	TqName       string `json:"tq_name" db:"tq_name"` 				//4	退签	货物退回发货人并签收
	TqTime       string `json:"tq_time" db:"tq_time"`
	PjName       string `json:"pj_name" db:"pj_name"` 				//5	派件	货物正在进行派件
	PjTime       string `json:"pj_time" db:"pj_time"`
	ThName       string `json:"th_name" db:"th_name"` 				//6	退回	货物正处于返回发货人的途中
	ThTime       string `json:"th_time" db:"th_time"`
	Created      string `json:"created" db:"created"`
	Updated      string `json:"updated" db:"updated"`
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
	Result     bool   `json:"result" db:"result"`
	ReturnCode string `json:"returnCode" db:"returnCode"`
	Message    string `json:"message" db:"message"`
}

// 对账存储结构
type KdSubscribeLog struct {
	LogisticsOrder string `db:"logistics_order"`
	Result     bool   `db:"result"`
	ReturnCode string `db:"returnCode"`
	Message    string `db:"message"`
}
