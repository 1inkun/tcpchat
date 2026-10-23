# TCP聊天室

## 基本信息

基于 `TCP` 连接的一个在线聊天室,有基础的限流和简单的账户系统,封禁名单由Redis维护

## 要求

- Go >= 1.26
- Make
- Redis

## 启动

### Redis

`/deployments` 文件夹下提供了运行在 `6789` 端口上密码为 `a123456` 的Redis配置文件

```bash
# 使用以下命令启动Redis
cd deployments
docker compose up -d
```

### 配置文件

项目中有硬编码的默认设置,如果需要自定义设置,请在项目根目录创建 `config.yaml` , 提供如下配置。如要使用默认值请不要创建 `config.yaml`:

```yaml
Server:
    Port: "服务器运行的端口,需要以':'为前缀"
    Safemode: "是否在log中显示连接用户的IP信息"
    TimeWindow: "运行用户两次消息间隔的时间,以秒为单位"
    MsgLimit: "从出现提示到封禁账户容许的消息条数"
    Timeout: "数据库查询的超时时间,以秒为单位"
    BanTime: "对用户的封禁时间,以分钟为单位"
Redis:
    Addr: "Redis服务器的IP地址,默认为'localhost:6379'"
    Password: "Redis的密码"
```

### 启动开发服务器

```bash
make dev
```

### 编译

```bash
make build
```

### 连接到服务器

```bash
telnet "你的聊天服务器IP" "你的端口"
```
