# banzhu-server

配置文件示例
```yaml
# 服务器配置
server:
  addr: ":8080"

log:
  level: "info"

# 数据库配置
database:
  dsn: "root:123456@tcp(192.168.10.199:3307)/book_db?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 10
  max_idle_conns: 5

# jwt令牌配置
jwt:
  access_secret: "bfdhbdhgdhd"
  access_expire_hours: 1
  refresh_secret: "bgfdhdsgrwtw"
  refresh_expire_hours: 720

# 文件存储配置
path:
  book: "./data/texts"
  chapter: "./data/texts/{}"
  images: "./data/images"

```