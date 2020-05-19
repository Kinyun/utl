package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vickydk/utl/config"
	"strings"
	"time"
)

func New() (*sql.DB, error) {
	db, err := sql.Open(config.Env.DBType, initDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(config.Env.MaxOpenConn)
	db.SetMaxIdleConns(config.Env.MaxIdle)
	err = db.Ping()
	if err != nil {
		fmt.Println("db is not connected")
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}

func initDSN() string {
	var psn strings.Builder
	psn.WriteString(config.Env.DBUser)
	psn.WriteString(":")
	psn.WriteString(config.Env.DBPass)
	if len(config.Env.DBSocket) > 0 {
		psn.WriteString("@unix(")
		psn.WriteString(config.Env.DBSocket)
		psn.WriteString(")/")
		psn.WriteString(config.Env.DBName)
		psn.WriteString("?loc=Local")
	} else {
		psn.WriteString("@tcp(")
		psn.WriteString(config.Env.DBHost)
		psn.WriteString(")/")
		psn.WriteString(config.Env.DBName)
		psn.WriteString("?parseTime=true")
	}
	fmt.Println(psn.String())
	return psn.String()
}
