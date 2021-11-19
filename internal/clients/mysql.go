package clients

import (
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var WriteDBCli *gorm.DB
var ReadDBCli *gorm.DB
var err error

func InitDBClients() {
	WriteDBCli, err = GetSqlDriver(config.GlobalConfig.WriteDB)
	if err != nil {
		panic(err)
	}
	ReadDBCli, err = GetSqlDriver(config.GlobalConfig.ReadDB)
	if err != nil {
		panic(err)
	}
}

func GetSqlDriver(dbConf config.DBConfig) (*gorm.DB, error) {
	var dbDialector = getDbDialector(dbConf)
	fmt.Println("user:%v, passwd:%v, ip:%v, port:%v", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port)
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	//_ = gormDb.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", func(d *gorm.DB) {
	//	d.Statement.RaiseErrorOnNotFound = false
	//})

	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetMaxIdleConns(dbConf.MaxIdleConns)
		rawDb.SetMaxOpenConns(dbConf.MaxOpenConns)
		return gormDb, nil
	}
}

func getDbDialector(conf config.DBConfig) gorm.Dialector {
	var dbDialector gorm.Dialector
	dsn := getDsn(conf)
	dbDialector = mysql.Open(dsn)
	return dbDialector
}

func getDsn(dbConf config.DBConfig) string {
	Host := dbConf.Host
	DataBase := dbConf.Name
	Port := dbConf.Port
	User := dbConf.User
	Pass := dbConf.Password
	Charset := "utf8mb4"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", User, Pass, Host, Port, DataBase, Charset)
}
