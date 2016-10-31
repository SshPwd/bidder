package db

import (
	"database/sql"
	"fmt"
	"golang_bidder/config"
	"golang_bidder/logs"
	_ "os"

	_ "github.com/go-sql-driver/mysql"
)

func InitMysql() {

	var (
		db_master *sql.DB
		err       error
	)

	conf := config.Get()

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		conf.MysqlMaster.Login, conf.MysqlMaster.Password,
		conf.MysqlMaster.Host, conf.MysqlMaster.Port,
		conf.MysqlMaster.DbName)

	db_master, err = sql.Open("mysql", connection)
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot connect to MySQL master server: %s:%s", conf.MysqlMaster.Host, conf.MysqlMaster.Port))
		// os.Exit(1)
		return
	}

	err = db_master.Ping()
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot connect to MySQL master server: %s:%s", conf.MysqlMaster.Host, conf.MysqlMaster.Port))
		// os.Exit(1)
		return
	}

	db_master.Close()
}
