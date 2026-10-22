package handler

import (
	"database/sql"
	"log"
	"net"

	"github.com/1inkun/tcpchat/internal/models"
	"github.com/1inkun/tcpchat/internal/utils"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	rdb *redis.Client
	db  *sql.DB
}

func NewAuthHandler(rdb *redis.Client, db *sql.DB) *AuthHandler {
	return &AuthHandler{
		rdb: rdb,
		db:  db,
	}
}

func (h *AuthHandler) Auth(conn net.Conn) (string, bool) {
	msg := `请完成授权操作!
输入!reg跳转注册
输入!log跳转登录
输入!ext断开连接`
	msg += "\n"
	_, err := conn.Write([]byte(msg))
	attackCount := 0
	if err != nil {
		return "", false
	}
	for {
		buffer := make([]byte, 32)
		_, err := conn.Read(buffer)
		if err != nil {
			return "", false
		}
		// input := buffer[0 : n-2]
		if string(buffer[0]) != "!" {
			msg := "输入的命令有误!\n"
			conn.Write([]byte(msg))
			attackCount++
		} else {
			cmd := string(buffer[1:4])
			log.Println(cmd)
			switch cmd {
			case "reg":
				return h.register(conn)
			case "log":
				return h.login(conn)
			case "ext":
				return "", false
			default:
				msg := "输入的命令有误!\n"
				conn.Write([]byte(msg))
				attackCount++
			}
		}
		if attackCount >= 5 {
			conn.Write([]byte("短时间内请求过多!\n"))
			return "", false
		}
	}
}

func (h *AuthHandler) login(conn net.Conn) (string, bool) {
	var LoginInfo models.LoginInfo
	buffer := make([]byte, 32)
	msg := "请输入登录信息\n"
	conn.Write([]byte(msg))

	// 获取登录信息
	msg = "请输入用户名:"
	conn.Write([]byte(msg))
	n, err := conn.Read(buffer)
	if err != nil {
		log.Panicf("在读取用户名输入时出现错误:%s", err)
		return "", false
	}
	LoginInfo.Username = string(buffer[0 : n-2])

	msg = "请输入密码:"
	conn.Write([]byte(msg))
	n, err = conn.Read(buffer)
	if err != nil {
		log.Panicf("在读取密码输入时出现错误:%s", err)
		return "", false
	}
	LoginInfo.Password = string(buffer[0 : n-2])

	// 登录逻辑
	// log.Println(LoginInfo)
	var passwordHash string
	query := "SELECT password FROM user WHERE username = ?;"
	row := h.db.QueryRow(query, LoginInfo.Username)
	if row.Err() != nil {
		log.Printf("在登录时查询数据库出现错误:%s", row.Err())
		return "", false
	}
	if err = row.Scan(&passwordHash); err != nil {
		log.Printf("在登录时查询数据库出现错误:%s", err)
		conn.Write([]byte("用户不存在\n"))
		return "", false
	}
	if err = utils.CompareHashAndPassword(passwordHash, LoginInfo.Password); err != nil {
		return "", false
	}
	// 登录成功返回用户名
	return LoginInfo.Username, true
}
func (h *AuthHandler) register(conn net.Conn) (string, bool) {
	var regInfo models.RegisterInfo
	buffer := make([]byte, 32)
	msg := "请输入注册信息\n"
	conn.Write([]byte(msg))
	// 获取注册信息
	for {
		msg := "请输入用户名:"
		conn.Write([]byte(msg))
		n, err := conn.Read(buffer)
		if err != nil {
			log.Panicf("在读取用户名输入时出现错误:%s", err)
			return "", false
		}
		regInfo.Username = string(buffer[0 : n-2])

		msg = "请输入密码:"
		conn.Write([]byte(msg))
		n, err = conn.Read(buffer)
		if err != nil {
			log.Panicf("在读取密码输入时出现错误:%s", err)
			return "", false
		}
		regInfo.Password = string(buffer[0 : n-2])

		msg = "请确认密码:"
		conn.Write([]byte(msg))
		n, err = conn.Read(buffer)
		if err != nil {
			log.Panicf("在读取密码输入时出现错误:%s", err)
			return "", false
		}
		check := string(buffer[0 : n-2])
		if check == regInfo.Password {
			break
		} else {
			msg := "输入信息有误,请重新注册\n"
			conn.Write([]byte(msg))
			return "", false
		}
	}
	// 根据信息进行注册操作
	// log.Println(regInfo)
	// 检查用户名是否已被注册
	var userId int64
	query := "SELECT id FROM user WHERE username = ?;"
	row := h.db.QueryRow(query, regInfo.Username)
	if row.Err() != nil {
		log.Printf("在查询数据库时出现错误:%s", row.Err())
		return "", false
	}
	err := row.Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("在查询数据库时出现错误:%s", err)
		return "", false
	}
	if userId > 0 {
		conn.Write([]byte("用户名已被占用\n"))
		return "", false
	}
	passwdHash, err := utils.GeneratePasswordHash(regInfo.Password)
	if err != nil {
		return "", false
	}
	query = "INSERT INTO `user` (`username`, `password`) VALUES (?,?);"
	_, err = h.db.Exec(query, regInfo.Username, passwdHash)
	if err != nil {
		return "", false
	}
	return regInfo.Username, true
}
