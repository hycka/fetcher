package bbc

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
)

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
	p := PostFactory("https://www.bbc.com/zhongwen/simp/world-55653976")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := setDate(p); err != nil {
		t.Errorf("test setPost err: %v", doc)
	}
	want := "2021-01-13T20:06:24.000Z"
	if p.Date != want {
		t.Errorf("got: %v, want: %v", p.Title, want)
	}
}

func TestSetTitle(t *testing.T) {
	p := PostFactory("https://www.bbc.com/zhongwen/simp/world-55653976")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := setTitle(p); err != nil {
		t.Errorf("test setPost err: %v", doc)
	}
	want := "中国俄罗斯海军进入印度洋挑战美国印太布局"
	if p.Title != want {
		t.Errorf("got: %v, want: %v", p.Title, want)
	}
}
func TestSetPost(t *testing.T) {
	p := PostFactory("https://www.bbc.com/zhongwen/simp/world-55653976")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetPost(p); err != nil {
		t.Errorf("test setPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}

func TestBBC(t *testing.T) {
	p := PostFactory("https://www.bbc.com/zhongwen/simp/world-55653976")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := bbc(p)
	fmt.Println(tc)
}
