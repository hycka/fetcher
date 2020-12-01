package boxun

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
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

func TestSetTitle(t *testing.T) {
	p := PostFactory("https://www.boxun.com/news/gb/taiwan/2020/07/202007101236.shtml")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetTitle(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	want := "香港新增42例确诊：8例输入性 5例来自哈萨克斯坦"
	if p.Title != want {
		t.Errorf("got: %v, want: %v", p.Title, want)
	}
}

func TestBoxun(t *testing.T) {
	p := PostFactory("https://www.dwnews.com/全球/60203378")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := Boxun(p)
	fmt.Println(tc)
}
