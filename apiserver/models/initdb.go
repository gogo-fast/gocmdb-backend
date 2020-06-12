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
	db_host := utils.GlobalConfig.Section("mysql").Key("host").String()
	db_port := utils.GlobalConfig.Section("mysql").Key("port").String()
	db_name := utils.GlobalConfig.Section("mysql").Key("db_name").String()
	db_user := utils.GlobalConfig.Section("mysql").Key("db_user").String()
	db_pass := utils.GlobalConfig.Section("mysql").Key("db_password").String()
	max_conn := utils.GlobalConfig.Section("mysql").Key("max_conn").MustInt(10)
	max_idle := utils.GlobalConfig.Section("mysql").Key("max_idle").MustInt(5)

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
