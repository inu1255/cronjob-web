package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	. "github.com/inu1255/cronjob"
)

func writeToFile(filename string, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func main() {
	fetcher := NewFetcher()
	resp, err := fetcher.GET("http://wap.sc.10086.cn/wap/verifyCode?t=0.22761388126741777", nil)
	if err != nil {
		log.Println(err)
		return
	}
	writeToFile("a.png", resp.Body)
	fmt.Println(".....")
	code, err := ImageCode("a.png")
	if err != nil {
		log.Println(err)
	}
	log.Println(code)
	form := make(url.Values)
	form.Add("loginType", "0")
	form.Add("loginNum", "18782071219")
	form.Add("loginPwd", "111")
	form.Add("verifyCode", "1423")
	form.Add("remFlag", "true")
	form.Add("loginInpage", "flow")
	resp, err = fetcher.POST("http://wap.sc.10086.cn/wap/member/login.json", strings.NewReader(form.Encode()), func(req *http.Request) {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	})
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
}
