package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var (
	PORT       = ":8080"
	MODE       = false
	TIMEWINDOW = 2
	MSGLIMIT   = 8
)

func safeRemoteAddr(conn net.Conn) string {
	if MODE {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

type MessageType int

const (
	NewConnection MessageType = iota + 1
	NewMessage
	Delete
)

type Message struct {
	MessageType MessageType
	Conn        net.Conn
	Text        string
}

type Client struct {
	Conn        net.Conn
	LastMessage time.Time
	StrikCount  int
}

func server(message chan Message) {
	// conns := map[string]net.Conn{}
	clients := map[string]Client{}
	for {
		msg := <-message
		switch msg.MessageType {
		case NewConnection:
			log.Printf("来自%s的新连接!", safeRemoteAddr(msg.Conn))
			// conns[msg.Conn.RemoteAddr().String()] = msg.Conn
			clients[msg.Conn.RemoteAddr().String()] = Client{
				Conn: msg.Conn,
			}
			// log.Println(conns)
		case NewMessage:
			// 消息限流
			now := time.Now()
			addr := msg.Conn.RemoteAddr().String()
			if now.Sub(clients[addr].LastMessage).Seconds() >= float64(TIMEWINDOW) {
				go func() {
					log.Printf("来自%s的新消息:%s", safeRemoteAddr(msg.Conn), msg.Text)
					// log.Println(msg.Text)
					newClient := Client{
						Conn:        msg.Conn,
						LastMessage: now,
						StrikCount:  0,
					}
					clients[addr] = newClient
					for key, client := range clients {
						log.Printf("向%s发送消息", key)
						if key == msg.Conn.RemoteAddr().String() {
							continue
						}
						n, err := client.Conn.Write([]byte(msg.Text))
						if err != nil {
							log.Printf("向客户端发送消息失败:%s", err)
							return
						}
						if n < len(msg.Text) {
							log.Println("向客户端发送消息时出现错误")
							return
						}
					}
				}()
			} else {
				// 当前时间距离上次发消息间隔不足两秒，直接拒绝
				clientInfo := clients[addr]
				if clientInfo.StrikCount%3 == 0 {
					msg := fmt.Sprintf("服务器规定在%d秒内只允许发送不超过%d条消息!", TIMEWINDOW, MSGLIMIT)
					_, err := clientInfo.Conn.Write([]byte(msg))
					if err != nil {
						log.Println("向客户端发送消息时出错")
						return
					}
				}
				clientInfo.StrikCount++
				if clientInfo.StrikCount == 0 || clientInfo.StrikCount >= 8 {
					// 攻击服务器大于两秒八次直接草飞
					clientInfo.Conn.Write([]byte("请停止攻击服务器!"))
					clientInfo.Conn.Close()
				}
				clients[addr] = clientInfo
			}
		case Delete:
			log.Printf("%s失去连接", safeRemoteAddr(msg.Conn))
			// delete(conns, msg.Conn.RemoteAddr().String())
			delete(clients, msg.Conn.RemoteAddr().String())
			msg.Conn.Close()
		}
	}
}

func welcome(conn net.Conn) {
	msg := "欢迎加入聊天室\n"
	remoteIP := safeRemoteAddr(conn)
	log.Printf("响应来自%s的连接", remoteIP)

	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("响应连接时失败:%s", err)
	}
	if n < len(msg) {
		log.Println("数据响应不完全")
	}
}

func client(conn net.Conn, messages chan Message) {
	welcome(conn)
	readBuff := make([]byte, 128)
	for {
		n, err := conn.Read(readBuff)
		if err != nil {
			log.Printf("客户端输入读取失败:%s", err)
			newMsg := Message{
				MessageType: Delete,
				Conn:        conn,
			}
			messages <- newMsg
			return
		} else {
			newMsg := Message{
				MessageType: NewMessage,
				Text:        string(readBuff[0:n]),
				Conn:        conn,
			}
			messages <- newMsg
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("服务器启动失败:%s", err)
	}
	log.Printf("服务器运行在端口%s上", PORT)

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("服务器连接处理失败:%s", err)
		}
		newMsg := Message{
			MessageType: NewConnection,
			Conn:        conn,
		}
		messages <- newMsg
		go client(conn, messages)
	}
}
