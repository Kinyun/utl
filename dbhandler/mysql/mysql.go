package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Optional struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
}

func New(opts *Optional) (*sql.DB, error) {
	db, err := sql.Open("mysql", opts.dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(opts.maxOpenConn)
	db.SetMaxIdleConns(opts.maxIdleConn)
	err = db.Ping()
	if err != nil {
		fmt.Println("db is not connected")
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}
