package utils

import (
	"net"

	"github.com/1inkun/tcpchat/config"
)

func SensMsg(conn net.Conn, cfg config.Config) string {
	// 处理敏感消息
	MODE := cfg.Server.SAFEMODE
	if MODE {
		return "[ 敏感信息 ]"
	} else {
		return conn.RemoteAddr().String()
	}
}
