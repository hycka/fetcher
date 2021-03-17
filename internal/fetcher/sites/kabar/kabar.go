package kabar

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
	p.Err = setDate(p)
	p.Err = setTitle(p)
	p.Err = setBody(p)
	return p.Err
}

func setDate(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	doc := htmldoc.ElementsByTagAndClass(p.DOC, "span", "article-date")
	d := []string{}
	if doc == nil {
		return fmt.Errorf("there is no element <time>")
	}
	//focus on node like "<span class="article-date"><i class="fa fa-clock-o"></i> 10/03/21 22:00 </span>"
	if doc[0].LastChild.Data != "" {
		d = append(d, doc[0].LastChild.Data)
	}
	//transform date to RFC3339 format
	t := strings.TrimSpace(d[0])
	tm, _ := time.Parse("02/01/06 15:04", t)
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
	b, err := kabar(p)
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

func kabar(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "post-content clearfix")
	if len(nodes) == 0 {
		nodes = htmldoc.ElementsByTagAndClass(doc, "div", "article-content-rawhtml")
	}
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
