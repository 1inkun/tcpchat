package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/1inkun/tcpchat/config"
	_ "github.com/mattn/go-sqlite3"
)

func InitDBConn(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Printf("连接数据库时出现错误:%s", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("数据库响应超时:%s", err)
		return nil, err
	}
	err = DBInit(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBInit(db *sql.DB) error {
	raw, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Printf("数据库schema文件读取失败:%s", err)
		return err
	}
	query := string(raw)
	// log.Println(query)
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("初始化数据库失败:%s", err)
		return err
	}
	return nil
}
