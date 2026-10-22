package handler

import (
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/models"
	"github.com/1inkun/tcpchat/internal/utils"
)

func Client(conn net.Conn, messages chan models.Message, cfg config.Config) {
	utils.Welcome(conn, cfg)
	readBuff := make([]byte, 128)
	for {
		n, err := conn.Read(readBuff)
		if err != nil {
			log.Printf("客户端输入读取失败:%s", err)
			newMsg := models.Message{
				MessageType: models.Delete,
				Conn:        conn,
			}
			messages <- newMsg
			return
		} else {
			newMsg := models.Message{
				MessageType: models.NewMessage,
				Text:        string(readBuff[0:n]),
				Conn:        conn,
			}
			messages <- newMsg
		}
	}
}
