package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gogo-cmdb/apiserver/utils"
	"os"
)

var db *sqlx.DB

func init() {
	var err error
	db_host := utils.GlobalConfig.GetString("mysql.host")
	db_port := utils.GlobalConfig.GetString("mysql.port")
	db_name := utils.GlobalConfig.GetString("mysql.db_name")
	db_user := utils.GlobalConfig.GetString("mysql.db_user")
	db_pass := utils.GlobalConfig.GetString("mysql.db_password")
	max_conn := utils.GlobalConfig.GetInt("mysql.max_conn")
	max_idle := utils.GlobalConfig.GetInt("mysql.max_idle")

	/*
		dsn format:
			"userName:password@tcp(dbHose:dbPort)/dbName?charset=utf8mb4&loc=Local&parseTime=true"
	*/
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		db_user, db_pass, db_host, db_port, db_name)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		utils.Logger.Error(err, "+")
		os.Exit(-1)
	}
	db.SetMaxIdleConns(max_conn)
	db.SetMaxIdleConns(max_idle)
	utils.Logger.Info("Initialize DB Success")

}
