package model

func InitModel() error {
	var err error
	//初始化Redis
	InitRedis()

	//初始SqlServer
	err = InitMsSql()
	return  err
}


type TradeOrder struct {
	Shelp          string `db:"shelp"`
	LogisticsComp  string `db:"logistics_comp"`
	LogisticsOrder string `db:"logistics_order"`
}
