package zaobao

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hi20160616/fetcher/internal/htmldoc"
	"github.com/hi20160616/gears"
	"golang.org/x/net/html"
)

type Post struct {
	Domain   string
	URL      *url.URL
	DOC      *html.Node
	Raw      []byte
	Title    string
	Body     string
	Date     string
	Filename string
}

func SetPost(p *Post) error {
	if err := SetDate(p); err != nil {
		return err
	}
	if err := SetTitle(p); err != nil {
		return err
	}
	if err := SetBody(p); err != nil {
		return err
	}
	return nil
}

func SetDate(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	metas := htmldoc.MetasByProperty(p.DOC, "article:modified_time")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("dwnews SetData got nothing.")
	}
	p.Date = cs[0]
	return nil
}

func SetTitle(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	n := htmldoc.ElementsByTag(p.DOC, "title")
	if n == nil {
		return fmt.Errorf("[-] there is no element <title>")
	}
	title := n[0].FirstChild.Data
	title = strings.ReplaceAll(title, " | 联合早报网", "")
	title = strings.ReplaceAll(title, " | 早报", "")
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	b, err := Zaobao(p)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = h1 + "\n\n" + b + "\n\n原地址：" + p.URL.String()
	return nil
}

func Zaobao(p *Post) (string, error) {
	if p.DOC == nil {
		return "", fmt.Errorf("[-] p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "article-content-container")
	if len(nodes) == 0 {
		nodes = htmldoc.ElementsByTagAndClass(doc, "div", "article-content-rawhtml")
	}
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<article>` from: " + p.URL.String())
	}
	plist := htmldoc.ElementsByTag(nodes[0], "p")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.FirstChild != nil &&
			v.FirstChild.Data == "strong" {
			a := htmldoc.ElementsByTag(v, "span")
			for _, aa := range a {
				body += aa.FirstChild.Data
			}
			body += "  \n"
		} else {
			body += v.FirstChild.Data + "  \n"
		}
	}
	body = strings.ReplaceAll(body, "span  \n", "")
	return body, nil
}
