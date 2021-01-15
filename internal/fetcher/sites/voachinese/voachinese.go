package voachinese

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
		return fmt.Errorf("p.DOC is nil")
	}
	doc := htmldoc.ElementsByTag(p.DOC, "time")
	// p.Date = doc[0].Attr[0].Val // short but not robust enough
	d := []string{}
	if doc == nil {
		return fmt.Errorf("there is no element <time>")
	}
	for _, a := range doc[0].Attr {
		if a.Key == "datetime" {
			d = append(d, a.Val)
		}
	}
	p.Date = d[0]
	return nil
}

func SetTitle(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	n := htmldoc.ElementsByTag(p.DOC, "title")
	if n == nil {
		return fmt.Errorf("there is no element <title>")
	}
	title := n[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	b, err := Voa(p)
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

func Voa(p *Post) (string, error) {
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")

	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "wsw")
	if nodes == nil {
		return "", errors.New(`There is no element match '<div class="wsw">'`)
	}
	plist := htmldoc.ElementsByTag(nodes[0], "p")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		}
		body += v.FirstChild.Data + "  \n"
	}
	body = strings.ReplaceAll(body, "strong  \n", "")
	body = strings.ReplaceAll(body, "span  \n", "")
	body = strings.ReplaceAll(body, "br  \n", "")
	return body, nil
}
