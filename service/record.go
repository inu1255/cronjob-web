package service

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/inu1255/gev2/models"
)

type Record struct {
	models.ModelExt `json:"-" xorm:"-"`
	Title           string            `json:"title,omitempty" xorm:"pk" gev:"名字"`
	UserId          int               `json:"user_id,omitempty" xorm:"pk" gev:"用户id"`
	UserParam       map[string]string `json:"params,omitempty" xorm:"" gev:"用户填写的参数"`
	Cookies         map[string]string `json:"cookies,omitempty" xorm:"" gev:"cookie"`
	Body            []byte            `json:"body,omitempty" xorm:"" gev:"返回值"`
	CreateAt        time.Time         `json:"create_at,omitempty" xorm:"created"`
	UpdateAt        time.Time         `json:"-" xorm:"updated"`
}

func (this *Record) ParseParam(param map[string]string) {
	if param != nil {
		if this.UserParam == nil {
			this.UserParam = param
		} else {
			for k, v := range param {
				this.UserParam[k] = v
			}
		}
	}
}

func (this *Record) ParseResponse(resp *http.Response) {
	if this.Cookies == nil {
		this.Cookies = make(map[string]string)
	}
	for _, item := range resp.Cookies() {
		this.Cookies[item.Name] = item.Value
	}
	this.Body, _ = ioutil.ReadAll(resp.Body)
}

func (this *Record) GetCookies() []*http.Cookie {
	if this.Cookies != nil {
		cookies := make([]*http.Cookie, 0, len(this.Cookies))
		for key, value := range this.Cookies {
			cookies = append(cookies, &http.Cookie{Name: key, Value: value, Path: "/", HttpOnly: false})
		}
		return cookies
	}
	return nil
}
