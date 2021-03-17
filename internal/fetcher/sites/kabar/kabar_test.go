package kabar

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/hi20160616/fetcher/internal/htmldoc"
)

var p = PostFactory("http://cn.kabar.kg/news/cn2021-03-10-05/")

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
	want := "2020-08-30T07:48:25+08:00"
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

func TestKabar(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := kabar(p)
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
