package main

import (
	"github.com/inu1255/cronjob/service"
	"github.com/inu1255/gev2"
	"github.com/inu1255/gev2/config"
	m "github.com/inu1255/gev2/models"
)

func main() {
	config.SetDb("mysql", "root:199337@/cronjob")

	verify := new(service.Verify)
	user := &service.User{}
	m.UserVerify = verify

	gev2.App.Static("/upload", "upload")
	gev2.App.Static("/api", "api")
	gev2.App.Use(m.UserMW(user))
	gev2.App.Use(gev2.CrossDomainMW())

	gev2.Bind("/user", user, "用户")
	gev2.Bind("/service", new(service.Service), "服务")
	gev2.Bind("/record", new(service.Record), "记录")
	gev2.Bind("/verify", verify, "验证码")

	gev2.Run(":8020")
}
