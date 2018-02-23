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

	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d",
		system.GetConfiguration().MsSqlServer,
		system.GetConfiguration().MsSqlDataBase,
		system.GetConfiguration().MsSqlUid,
		system.GetConfiguration().MsSqlPwd,
		system.GetConfiguration().MsSqlPort)
	fmt.Println(connString)
	conn, err := sqlx.Connect("mssql", connString)
	if err != nil {
		return errors.New(fmt.Sprintf("Open connection failed:", err.Error()))
	}
	DB = conn
	fmt.Println(DB)
	return err
}


func (tradeOrder *TradeOrder) GetTopOrder() ([]TradeOrder, error) {

	rows, err := DB.Queryx("SELECT top 10 shelp,logistics_comp,logistics_order FROM trade_order;")
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
		}
		result = append(result, row)
	}
	return result, err
}
