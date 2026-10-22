package utils

import (
	"log"
	"net"

	"github.com/1inkun/tcpchat/config"
)

func Welcome(conn net.Conn, cfg config.Config) {
	msg := ` _ _ _       _
| | | | ___ | | ___  ___ ._ _ _  ___
| | | |/ ._>| |/ | '/ . \| ' ' |/ ._>
|__/_/ \___.|_|\_|_.\___/|_|_|_|\___.
                                     `
	// remoteIP := safeRemoteAddr(conn)
	msg = msg + "\n"
	log.Printf("响应来自%s的连接", SensMsg(conn, cfg))

	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("响应连接时失败:%s", err)
	}
	if n < len(msg) {
		log.Println("数据响应不完全")
	}
}
