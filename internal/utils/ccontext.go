package utils

import (
	"context"
	"net"
	"time"

	"github.com/1inkun/tcpchat/config"
)

type ChatContext struct {
	Conn    net.Conn
	Context context.Context
}

func NewConn(conn net.Conn, cfg config.Config) *ChatContext {
	c := context.Background()
	timeout := time.Duration(cfg.Server.TIMEOUT) * time.Second
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	return &ChatContext{
		Conn:    conn,
		Context: ctx,
	}
}
