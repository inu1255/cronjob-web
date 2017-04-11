package service

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/inu1255/gev2/models"
)

type IService interface {
}

type Service struct {
	models.ModelExt `json:"-" xorm:"-"`
	Title           string            `json:"title,omitempty" xorm:"pk" gev:"名字"`
	Method          string            `json:"method,omitempty" xorm:"varchar(32)" gev:"GET/POST"`
	Url             string            `json:"url,omitempty" xorm:"" gev:"链接"`
	UserParam       map[string]string `json:"user_param,omitempty" xorm:"" gev:"用户参数"`
	Params          string            `json:"params,omitempty" xorm:"text" gev:"参数"`
	Header          map[string]string `json:"header,omitempty" xorm:"" gev:"头部"`
	ValidRegex      string            `json:"valid_regex,omitempty" xorm:"" gev:"判断是否成功的正则表达式"`
	Script          string            `json:"script,omitempty" xorm:"" gev:"TODO: 用于生成其它参数的脚本"`
	Parent          string            `json:"parent,omitempty" xorm:"" gev:"依赖的service的id"`
	Type            int               `json:"type,omitempty" xorm:"" gev:"类型 0:辅助 1:项目"`

	record *Record `json:"-" xorm:"-" gev:"参数及返回值"`
}

// 生成 http 请求的 body
func (this *Service) ParseBody(param map[string]string) io.Reader {
	s := this.Params
	if param != nil {
		for k, v := range param {
			s = strings.Replace(s, k, v, -1)
		}
	}
	return strings.NewReader(s)
}

// 获取父Service
func (this *Service) GetParent() (*Service, bool) {
	if this.Parent != "" {
		service := new(Service)
		ok, err := Rdb().Where("title=?", this.Parent).Get(service)
		if err != nil {
			Log.Println(err)
		}
		if ok {
			this.NewModel(service)
			return service, true
		}
	}
	return nil, false
}

func (this *Service) GetCookies(user_id int) ([]*http.Cookie, bool) {
	if this.Parent != "" && user_id > 0 {
		record := new(Record)
		ok, err := Rdb().Where("title=? and user_id=?", this.Parent, user_id).Get(record)
		if err != nil {
			Log.Println(err)
		}
		if ok {
			return record.GetCookies(), true
		}
	}
	return nil, false
}

func (this *Service) IsValid(body []byte) bool {
	if this.ValidRegex == "" {
		return true
	}
	reg, err := regexp.Compile(this.ValidRegex)
	if err != nil {
		return false
	}
	return reg.Match(body)
}

func (this *Service) GetRecord(user_id int) (*Record, bool) {
	if this.record != nil && this.record.UserId == user_id {
		return this.record, true
	}
	record := new(Record)
	record.Title = this.Title
	record.UserId = user_id
	ok, _ := this.Db.Get(record)
	this.record = record
	this.NewModel(record)
	return record, ok
}

// @path
func (this *Service) Run(user_id int, isFirst bool) (*Record, error) {
	client := http.DefaultClient
	record, _ := this.GetRecord(user_id)
	if this.UserParam != nil && len(this.UserParam) > 0 {
		for key, _ := range this.UserParam {
			if _, ok := record.UserParam[key]; !ok {
				return nil, models.ApiErr(this.Title+"::缺少参数::"+key, 1)
			}
		}
	}
	req, err := http.NewRequest(this.Method, this.Url, this.ParseBody(record.UserParam))
	if err != nil {
		Log.Println(err)
		return nil, models.ApiErr(this.Title+"::构造请求失败", 2)
	}
	if this.Header != nil {
		for key, value := range this.Header {
			req.Header.Set(key, value)
		}
	}
	if cookies, ok := this.GetCookies(user_id); ok {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		Log.Println(err)
		return nil, models.ApiErr(this.Title+"::请求失败", 3)
	}
	record.ParseResponse(resp)
	if !this.IsValid(record.Body) {
		if this.Parent != "" && isFirst {
			if parent, ok := this.GetParent(); ok {
				_, err = parent.Run(user_id, true)
				if err != nil {
					return nil, err
				}
				return this.Run(user_id, false)
			}
			return nil, models.ApiErr(this.Parent+"::不存在", 4)
		}
		return nil, models.ApiErr(this.Title+"::invalid", 5)
	}
	record.Save()
	return record, nil
}

// @param title service.title path
func (this *Service) Start(title string, param map[string]string) (*Record, error) {
	if user, ok := this.GetUser(); ok {
		if ok, err := this.Db.Where("title=?", title).Get(this); !ok {
			return nil, err
		}
		user_id := user.GetId()
		record, _ := this.GetRecord(user_id)
		record.ParseParam(param)
		record.Save()
		return this.Run(user_id, true)
	}
	return nil, models.NeedAuthError
}
