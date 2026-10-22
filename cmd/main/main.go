package main

import (
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/db"
	"github.com/1inkun/tcpchat/internal/handler"
	"github.com/1inkun/tcpchat/internal/models"
)

func main() {
	cfg := config.InitConfig()
	// log.Println(cfg)

	// 初始化数据库连接
	sqlDb, err := db.InitDBConn(*cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败:%s", err)
		return
	}
	defer sqlDb.Close()
	rdb, err := db.InitRedisConn(*cfg)
	if err != nil {
		log.Fatalf("redis初始化失败:%s", err)
		return
	}

	ln, err := net.Listen("tcp", cfg.Server.PORT)
	if err != nil {
		log.Fatalf("服务器启动失败:%s", err)
		return
	}
	log.Printf("服务器运行在端口%s上", cfg.Server.PORT)

	// 注册handler
	authHandler := handler.NewAuthHandler(rdb, sqlDb)

	messages := make(chan models.Message)
	go handler.Server(messages, *cfg)
	// 处理连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("服务器连接处理失败:%s", err)
		}
		go handler.Client(conn, messages, *cfg, *authHandler)
	}
}
