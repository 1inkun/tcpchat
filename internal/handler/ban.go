package handler

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/1inkun/tcpchat/config"
	"github.com/redis/go-redis/v9"
)

type BanHandler struct {
	rdb *redis.Client
	db  *sql.DB
}

func NewBanHandler(rdb *redis.Client, db *sql.DB) *BanHandler {
	return &BanHandler{
		rdb: rdb,
		db:  db,
	}
}

var layout = "2006-01-02 15:04:05"

func (h *BanHandler) Check(conn net.Conn, cfg config.Config) (string, bool) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Server.TIMEOUT))
	defer cancel()
	remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	// log.Printf("查询封禁情况的IP地址:%s", remoteIP)
	res, err := h.rdb.Get(ctx, remoteIP).Result()
	if err != nil {
		if err == redis.Nil {
			return "", true
		}
		log.Printf("在查询封禁列表时遇到错误:%s", err)
		return "", false
	}

	// log.Printf("在redis中查询到的封禁信息:%s", res)
	releaseTime, err := time.Parse("", res)
	if releaseTime.After(now) {
		// 封禁时限已到,解除封禁
		// log.Println("封禁时限已到,解除封禁")
		_, err := h.rdb.Del(ctx, remoteIP).Result()
		if err != nil {
			log.Printf("在处理封禁列表时遇到错误:%s", err)
			return "", true
		}
		return "", true
	} else {
		// 时限未到,返回封禁时间
		// log.Printf("封禁时限未到,当前时间:%s", now.Format(layout))
		return releaseTime.Format(layout), false
	}
}
func (h *BanHandler) GetBaned(conn net.Conn, cfg config.Config) (string, bool) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Server.TIMEOUT))
	defer cancel()
	exp := time.Minute * time.Duration(cfg.Server.BANTIME)
	banTime := now.Add(exp).String()
	remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	// log.Printf("被封禁的IP地址:%s", remoteIP)
	// 向redis中添加封禁信息
	_, err := h.rdb.Set(ctx, remoteIP, banTime, exp).Result()
	if err != nil {
		log.Printf("在处理封禁列表时遇到错误:%s", err)
		return "", false
	}
	return now.Add(exp).Format(layout), true
}
