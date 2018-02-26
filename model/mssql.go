package model

import (
	//"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"kuaidi100_go/system"
)

var DB *sqlx.DB

func InitMsSql() error {

	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;",
		system.GetConfiguration().MsSqlServer,
		system.GetConfiguration().MsSqlDataBase,
		system.GetConfiguration().MsSqlUid,
		system.GetConfiguration().MsSqlPwd,
		system.GetConfiguration().MsSqlPort)
	conn, err := sqlx.Connect("mssql", connString)
	fmt.Println(connString)
	if err != nil {
		return errors.New(fmt.Sprintf("Open connection failed:", err.Error()))
	}
	DB = conn
	return err
}


func (tradeOrder *TradeOrder) GetTopOrder() ([]TradeOrder, error) {
	//start_date:="where create_date >= '2018-01-01 00:00:00'  "
	//rows, err := DB.Queryx("SELECT * FROM trade_order WHERE site_order_id IN ('412151205240974365','412329093254557643');")
	//rows, err := DB.Queryx("SELECT top 5 shelp,logistics_comp,logistics_order FROM trade_order; ")
	rows, err := DB.Queryx("SELECT top 5 shelp,logistics_comp,logistics_order FROM trade_order where shelp='DISTRIBUTOR_13174102'; ")
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
