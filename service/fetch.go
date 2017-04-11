package service

import (
	"io"
	"net/http"
)

type Fetcher struct {
	cookies []*http.Cookie
	client  *http.Client
}

func NewFetcher() *Fetcher {
	this := &Fetcher{}
	this.client = &http.Client{}
	return this
}

func (this *Fetcher) Fetch(method, url string, body io.Reader, call func(*http.Request)) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if this.cookies != nil {
		for _, c := range this.cookies {
			req.AddCookie(c)
		}
	}
	if call != nil {
		call(req)
	}
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	this.cookies = resp.Cookies()
	return resp, nil
}

func (this *Fetcher) GET(url string, call func(*http.Request)) (*http.Response, error) {
	return this.Fetch("GET", url, nil, call)
}

func (this *Fetcher) POST(url string, body io.Reader, call func(*http.Request)) (*http.Response, error) {
	return this.Fetch("POST", url, body, call)
}
