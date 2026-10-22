package main

import (
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/handler"
	"github.com/1inkun/tcpchat/internal/models"
)

func main() {
	cfg := config.InitConfig()
	// log.Println(cfg)

	ln, err := net.Listen("tcp", cfg.Server.PORT)
	if err != nil {
		log.Fatalf("服务器启动失败:%s", err)
	}
	log.Printf("服务器运行在端口%s上", cfg.Server.PORT)

	messages := make(chan models.Message)
	go handler.Server(messages, *cfg)
	// 处理连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("服务器连接处理失败:%s", err)
		}
		newMsg := models.Message{
			MessageType: models.NewConnection,
			Conn:        conn,
		}
		messages <- newMsg
		go handler.Client(conn, messages, *cfg)
	}
}
