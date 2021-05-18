package cna

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
)

var p = PostFactory("https://www.cna.com.tw/news/aopl/202009300058.aspx")

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
	if err := setDate(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	want := "2020-09-30T10:54:00+08:00"
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
		t.Errorf("test SetPost err: %v", err)
	}
	want := "被爆10年沒繳稅 川普：避稅計畫展現我的才智 | 國際"
	if p.Title != want {
		t.Errorf("\ngot: %v\nwant: %v", p.Title, want)
	}
}

func TestCna(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := cna(p)
	fmt.Println(tc)
}

func TestSetPost(t *testing.T) {
	// var p = PostFactory("https://www.cna.com.tw/news/afe/202009290241.aspx") // should be ignore
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetPost(p); err != nil {
		t.Errorf("test SetPost err: %v", err)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}

func TestTranform(t *testing.T) {
	p := PostFactory("https://www.cna.com.tw/news/afe/202009290241.aspx")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err = SetPost(p); err != nil {
		t.Error(err)
	}
	if err := transform(p); err != nil {
		t.Errorf("test SetPost err: %v", err)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}
