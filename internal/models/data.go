package models

import "time"

type User struct {
	Id            int64     `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	LastLoginIP   string    `json:"last_login_ip"`
	LastLoginDate time.Time `json:"last_login_date"`
}

type News struct {
	Id         int64     `json:"id"`
	Content    string    `json:"content"`
	Updatetime time.Time `json:"update_time"`
}

type RegisterInfo struct {
	Username string
	Password string
}

type LoginInfo struct {
	Username string
	Password string
}
