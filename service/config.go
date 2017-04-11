package service

import (
	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev2/config"
)

var (
	Db  *xorm.Engine
	Log = config.NewLogger("cronjob")
)

func Rdb() *xorm.Session {
	return config.Rdb()
}
