package reuters

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
)

var p = PostFactory("https://cn.reuters.com/theWire//article/china-bond-market-eu-pension-fund-0120-idCNKBS29P0XU")

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
		t.Errorf("test setPost err: %v", doc)
	}
	want := "2021-01-20T09:15:12+0000"
	if p.Date != want {
		t.Errorf("got: %v, want: %v", p.Date, want)
	}
}

func TestSetTitle(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := setTitle(p); err != nil {
		t.Errorf("test setPost err: %v", doc)
	}
	want := "焦点：中国空前减税降费 企业为何感受却不强"
	if p.Title != want {
		t.Errorf("got: %v, want: %v", p.Title, want)
	}
}
func TestSetPost(t *testing.T) {
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

func TestReuters(t *testing.T) {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := reuters(p)
	fmt.Println(tc)
}
