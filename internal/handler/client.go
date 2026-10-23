package handler

import (
	"fmt"
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/models"
	"github.com/1inkun/tcpchat/internal/utils"
)

func Client(conn net.Conn, messages chan models.Message, cfg config.Config, auth AuthHandler, ban BanHandler) {
	// 先检查当前IP是否被封禁
	if bandTime, ok := ban.Check(conn, cfg); !ok {
		// 未通过授权直接断开连接
		log.Printf("%s 已被封禁", utils.SensMsg(conn, cfg))
		resp := fmt.Sprintf("当前IP已被封禁! 将于%s解封!\n", bandTime)
		conn.Write([]byte(resp))
		conn.Close()
		return
	}
	// 检查授权状态
	username, ok := auth.Auth(conn, ban, cfg)
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
