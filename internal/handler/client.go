package handler

import (
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/models"
	"github.com/1inkun/tcpchat/internal/utils"
)

func Client(conn net.Conn, messages chan models.Message, cfg config.Config, auth AuthHandler) {
	// 先检查授权状态
	username, ok := auth.Auth(conn)
	if !ok {
		// 未通过授权直接断开连接
		log.Printf("%s 未通过授权", utils.SensMsg(conn, cfg))
		conn.Close()
		return
	} else {
		newMsg := models.Message{
			MessageType: models.NewConnection,
			Conn:        conn,
			Username:    username,
		}
		messages <- newMsg
	}
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
				Username:    username,
			}
			messages <- newMsg
		}
	}
}
