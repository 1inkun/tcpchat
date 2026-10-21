package main

import (
	"log"
	"net"
)

var (
	PORT = ":8080"
	MODE = false
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

func server(message chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-message
		switch msg.MessageType {
		case NewConnection:
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
			// log.Println(conns)
		case NewMessage:
			go func() {
				// log.Println(msg.Text)
				for key, conn := range conns {
					if key == msg.Conn.RemoteAddr().String() {
						continue
					}
					n, err := conn.Write([]byte(msg.Text))
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
		case Delete:
			delete(conns, msg.Conn.RemoteAddr().String())
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
	readBuff := make([]byte, 4096)
	for {
		n, err := conn.Read(readBuff)
		if err != nil {
			log.Printf("客户端输入读取失败:%s", err)
			newMsg := Message{
				MessageType: Delete,
				Conn:        conn,
			}
			messages <- newMsg
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
