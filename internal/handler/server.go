package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/1inkun/tcpchat/config"
	"github.com/1inkun/tcpchat/internal/models"
	"github.com/1inkun/tcpchat/internal/utils"
)

func Server(message chan models.Message, ban BanHandler, cfg config.Config) {
	// conns := map[string]net.Conn{}
	clients := map[string]models.Client{}
	for {
		msg := <-message
		switch msg.MessageType {
		case models.NewConnection:
			log.Printf("来自%s的新连接!", utils.SensMsg(msg.Conn, cfg))
			// conns[msg.Conn.RemoteAddr().String()] = msg.Conn
			// 同名用户登录挤下线
			for k, v := range clients {
				if v.Username == msg.Username {
					v.Conn.Write([]byte("当前账户已在其他地方登录\n"))
					v.Conn.Close()
					// 挤下线后在clients列表中删除
					delete(clients, k)
					break
				}
			}
			clients[msg.Conn.RemoteAddr().String()] = models.Client{
				Conn:     msg.Conn,
				Username: msg.Username,
			}
			utils.Welcome(msg.Conn, cfg)
			text := fmt.Sprintf("当前在线人数: %d 人\n", len(clients))
			msg.Conn.Write([]byte(text))
			// TODO: 新用户同步聊天记录
		case models.NewMessage:
			// 消息限流
			now := time.Now()
			addr := msg.Conn.RemoteAddr().String()
			if now.Sub(clients[addr].LastMessage).Seconds() >= float64(cfg.Server.TIMEWINDOW) {
				// 成功发送消息
				log.Printf("来自%s的新消息:%s", utils.SensMsg(msg.Conn, cfg), msg.Text)
				// 更新发送消息用户的信息
				newClient := models.Client{
					Conn:        msg.Conn,
					LastMessage: now,
					StrikCount:  0,
					Username:    msg.Username,
				}
				clients[addr] = newClient
				// 格式化消息
				text := fmt.Sprintf("%s %s: %s\n", now.Format("2006-01-02 15:04:05"), msg.Username, msg.Text)
				// 使用协程遍历当前已连接用户列表
				go func() {
					for key, client := range clients {
						// log.Printf("向%s发送消息", key)
						if key == msg.Conn.RemoteAddr().String() {
							continue
						}
						n, err := client.Conn.Write([]byte(text))
						if err != nil {
							log.Printf("向客户端发送消息失败:%s", err)
							break
						}
						if n < len(msg.Text) {
							log.Println("向客户端发送消息时出现错误")
							break
						}
					}
				}()
			} else {
				// 当前时间距离上次发消息间隔不足两秒，直接拒绝
				clientInfo := clients[addr]

				msg := fmt.Sprintf("服务器规定在%d秒内只允许发送不超过%d条消息!", cfg.Server.TIMEWINDOW, cfg.Server.MSGLIMIT)
				_, err := clientInfo.Conn.Write([]byte(msg))
				if err != nil {
					log.Println("向客户端发送消息时出错")
					continue
				}

				// TODO: 对多次爆破服务器的IP实施封禁,封禁名单由Redis维护
				clientInfo.StrikCount++
				if clientInfo.StrikCount >= cfg.Server.MSGLIMIT {
					// 攻击服务器大于两秒八次直接草飞
					if banTime, ok := ban.GetBaned(clientInfo.Conn, cfg); ok {
						resp := fmt.Sprintf("你已被封禁!将于%s解封\n", banTime)
						clientInfo.Conn.Write([]byte(resp))
						clientInfo.Conn.Close()
					}
					clientInfo.Conn.Write([]byte("请停止攻击服务器!\n"))
					clientInfo.Conn.Close()
				}
				clients[addr] = clientInfo
			}
		case models.Delete:
			log.Printf("%s失去连接", utils.SensMsg(msg.Conn, cfg))
			// delete(conns, msg.Conn.RemoteAddr().String())
			delete(clients, msg.Conn.RemoteAddr().String())
			msg.Conn.Close()
		}
	}
}
