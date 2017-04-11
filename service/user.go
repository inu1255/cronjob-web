package service

import (
	"github.com/inu1255/gev2/models"
)

type User struct {
	models.UserRegistModel `xorm:"extends"`
}

func (this *User) CopyTo(user models.IUserModel, bean interface{}) error {
	data := bean.(*User)
	data.Id = this.Id
	data.Nickname = this.Nickname
	return nil
}
