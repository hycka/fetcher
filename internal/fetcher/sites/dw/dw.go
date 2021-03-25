package dw

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
	Err      error
}

func SetPost(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	p.Err = setTitle(p)
	p.Err = setDate(p)
	p.Err = setBody(p)
	return p.Err
}

func setDate(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.Title == "" {
		return fmt.Errorf("p.Title is nil")
	}
	//focus on title like "港澳煞停接种德国BioNTech疫苗 | 德国之声 来自德国 介绍德国 | DW | 24.03.2021"
	s := strings.Split(p.Title, "｜")
	tmp := strings.TrimSpace(s[len(s)-1])
	//transform date to RFC3339 format   "2006-01-02 15:04:05"
	tm, err := time.Parse("02.01.2006", tmp)
	if err != nil {
		return fmt.Errorf("can not get Date")
	}
	p.Date = tm.Format(time.RFC3339)
	return nil
}

func setTitle(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	doc := htmldoc.ElementsByTag(p.DOC, "title")
	if doc == nil {
		return fmt.Errorf("there is no element <title>")
	}
	title := doc[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func setBody(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	b, err := dw(p)
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

func dw(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "longText")
	if len(nodes) == 0 {
		return "", errors.New("There is no tag named `<article>` from: " + p.URL.String())
	}
	plist := htmldoc.ElementsByTag(nodes[0], "p")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		} else {
			body += v.FirstChild.Data + "  \n"
		}
	}
	// body = strings.ReplaceAll(body, "span  \n", "")
	return body, nil
}
