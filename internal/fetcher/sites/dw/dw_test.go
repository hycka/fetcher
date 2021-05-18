package dw

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
)

var p = PostFactory("https://www.dw.com/zh/%E9%A6%99%E6%B8%AF%E5%B0%8F%E5%AD%A6%E6%95%99%E6%9D%90%E5%88%A0%E9%99%A4%E4%B8%AD%E5%8D%8E%E6%B0%91%E5%9B%BD%E6%94%BF%E5%BA%9C/a-57168239")

// var p = PostFactory("https://www.dw.com/zh/%E6%AC%A7%E6%B4%B2%E4%BA%BA%E6%9D%83%E6%B3%95%E9%99%A2%E5%BC%BA%E5%88%B6%E6%8E%A5%E7%A7%8D%E7%96%AB%E8%8B%97%E4%B8%8D%E8%BF%9D%E5%8F%8D%E4%BA%BA%E6%9D%83/a-57144574")

// var p = PostFactory("https://www.dw.com/zh/%E6%8B%9C%E7%99%BB%E5%85%A8%E4%BD%93%E7%BE%8E%E5%9B%BD%E5%85%AC%E6%B0%91%E6%9C%AC%E6%9C%88%E5%B0%B1%E9%83%BD%E8%83%BD%E6%8E%A5%E7%A7%8D%E7%96%AB%E8%8B%97/a-57119062")

// var p = PostFactory("https://www.dw.com/zh/%E6%B6%89%E5%85%83%E6%9C%97%E4%BA%8B%E4%BB%B6%E8%B0%83%E6%9F%A5%E6%8A%A5%E5%AF%BC%E6%B8%AF%E5%8F%B0%E7%BC%96%E5%AF%BC%E8%99%9A%E5%81%87%E9%99%88%E8%BF%B0%E6%A1%88%E5%BC%80%E5%BA%AD/a-56971577")

func PostFactory(rawurl string) *Post {
	url, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Post{
		Domain: url.Hostname(),
		URL:    url,
	}
}

func TestSetDate(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	setTitle(p)
	if err := setDate(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	want := "2021-04-07T00:00:00Z"
	if p.Date != want {
		t.Errorf("\ngot: %v\nwant: %v", p.Date, want)
	}
}

func TestSetTitle(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := setTitle(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	want := "吉外交部长：任何国家和地区都无法独自对抗外部威胁"
	if p.Title != want {
		t.Errorf("\ngot: %v\nwant: %v", p.Title, want)
	}
}

func TestDw(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := dw(p)
	fmt.Println(tc)
}

func TestSetPost(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetPost(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}
