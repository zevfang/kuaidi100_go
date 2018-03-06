package model

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"fmt"
	"kuaidi100_go/system"
	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sqlx.DB

func InitMsSql() error {

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;",
		system.GetConfiguration().MsSqlServer,
		system.GetConfiguration().MsSqlPort,
		system.GetConfiguration().MsSqlDataBase,
		system.GetConfiguration().MsSqlUid,
		system.GetConfiguration().MsSqlPwd)

	conn, err := sqlx.Connect("mssql", connString)

	if err != nil {
		return errors.New(fmt.Sprintf("Open connection failed:", err.Error()))
	}
	DB = conn
	return err
}

func GetTopOrder() ([]TradeOrder, error) {

	top_count := system.GetConfiguration().OncePollCount
	if top_count == 0 {
		top_count = 100
	}

	start_data := system.GetConfiguration().StartDate
	if start_data == "" {
		start_data = "2018-01-01 00:00:00"
	}
	sql := fmt.Sprintf("select top %d shelp,logistics_comp,logistics_order FROM trade_order(nolock)  where is_subscribe = 0 and create_date >= ?1;", top_count)
	rows, err := DB.Queryx(sql, top_count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []TradeOrder
	for rows.Next() {
		var row TradeOrder
		err := rows.StructScan(&row)
		if err != nil {
			fmt.Println("row err")
			continue
		}
		result = append(result, row)
	}
	return result, err
}

func TranOrderAndResult(resultData KdSubscribeLog) error {
	tx, err := DB.Beginx()
	//更新订阅状态
	_, err = tx.Exec("UPDATE trade_order SET  is_subscribe = 1,subscribe_date = ?1 WHERE logistics_order=?2;", system.GetNow(), resultData.LogisticsOrder)
	//插入对账数据
	_, err = tx.Exec("INSERT INTO  kd_subscribe_log (logistics_order,result,returnCode,message,created) VALUES (?1,?2,?3,?4,?5);",
		resultData.LogisticsOrder, resultData.Result, resultData.ReturnCode, resultData.Message, system.GetNow())
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// 查询推送信息--根据快递单号记录
func GetKdOrderOne(nu string) (KdOrder, error) {
	kdOrder := KdOrder{}
	const sql = "SELECT * FROM  kd_order(NOLOCK) WHERE nu=?1"
	err := DB.Get(&kdOrder, sql, nu)
	return kdOrder, err
}

// 获取阿芙物流编码
func GetTradeOrderShelp(nu string) (string, error) {
	var shelp string
	const sql = ` SELECT top 1 shelp FROM trade_order(NOLOCK) WHERE logistics_order=?1`
	err := DB.QueryRow(sql, nu).Scan(&shelp)
	return shelp, err
}

// 新增推送信息
func InsertKdOrder(order KdOrder) error {
	const sql = `INSERT INTO [dbo].[kd_order] ([kd_status],[kd_status_name],[kd_message],[state],[state_name],[com],[shelp],[nu],
					[data],[zt_name],[zt_time],[lj_name],[lj_time],[yn_name],[yn_time],[qs_name],[qs_time],[tq_name],[tq_time],[pj_name],[pj_time],[th_name],[th_time],
					[created],[updated]) 
				VALUES	(:kd_status,:kd_status_name,:kd_message,:state,:state_name,:com,:shelp,:nu,:data,
					:zt_name,:zt_time,:lj_name,:lj_time,:yn_name,:yn_time,:qs_name,:qs_time,:tq_name,:tq_time,:pj_name,:pj_time,:th_name,:th_time,
					:created,:updated);`

	_, err := DB.NamedExec(sql, &order)
	return err
}

// 更新推送信息
func UpdateKdOrder(order KdOrder) error {
	const sql = ` UPDATE [dbo].[kd_order]
				   SET [kd_status] = :kd_status
					  ,[kd_status_name] = :kd_status_name
					  ,[kd_message] = :kd_message
					  ,[state] = :state
					  ,[state_name] = :state_name
					  ,[com] = :com
					  ,[shelp] = :shelp
					  ,[nu] = :nu
					  ,[data] = :data
					  ,[zt_name] = :zt_name
					  ,[zt_time] = :zt_time
					  ,[lj_name] = :lj_name
					  ,[lj_time] = :lj_time
					  ,[yn_name] = :yn_name
					  ,[yn_time] = :yn_time
					  ,[qs_name] = :qs_name
					  ,[qs_time] = :qs_time
					  ,[tq_name] = :tq_name
					  ,[tq_time] = :tq_time
					  ,[pj_name] = :pj_name
					  ,[pj_time] = :pj_time
					  ,[th_name] = :th_name
					  ,[th_time] = :th_time
					  ,[created] = :created
					  ,[updated] = :updated
  					WHERE [nu] = :nu;`
	_, err := DB.NamedExec(sql, &order)
	return err
}
