package dal

import (
	"clicky.website/clicky/security/biz/dal/mysql"
	"clicky.website/clicky/security/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
