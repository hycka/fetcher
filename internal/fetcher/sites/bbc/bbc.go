package bbc

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
	metas := htmldoc.MetasByName(p.DOC, "article:modified_time")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("bbc SetData got nothing.")
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
	if strings.Contains(title, "[图集]") {
		return fmt.Errorf("[!] Picture news ignored.")
	}
	title = strings.ReplaceAll(title, " - BBC News 中文", "")
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	b, err := BBC(p)
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

func BBC(p *Post) (string, error) {
	if p.DOC == nil {
		return "", fmt.Errorf("[-] p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTag(doc, "main")
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<article>` from: " + p.URL.String())
	}
	articleDoc := nodes[0]
	plist := htmldoc.ElementsByTag(articleDoc, "h2", "p")

	for _, v := range plist {
		if v.FirstChild != nil {
			if v.Parent.FirstChild.Data == "h2" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.Data)
			} else if v.FirstChild.Data == "b" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.FirstChild.Data)
			} else {
				body += v.FirstChild.Data + "  \n"
			}
		}
	}
	body = strings.ReplaceAll(body, "span  \n", "")
	return body, nil
}
