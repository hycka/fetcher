package dwnews

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
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
	metas := htmldoc.MetasByName(p.DOC, "parsely-pub-date")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("dwnews setData got nothing.")
	}
	p.Date = cs[0]
	return nil
}

func setTitle(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	n := htmldoc.ElementsByTag(p.DOC, "title")
	if n == nil {
		return fmt.Errorf("there is no element <title>")
	}
	title := n[0].FirstChild.Data
	if strings.Contains(title, "[图集]") {
		return fmt.Errorf("[!] Picture news ignored.")
	}
	title = strings.ReplaceAll(title, "｜多维新闻", "")
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
	b, err := dwnews(p)
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

func dwnews(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTag(doc, "article")
	if len(nodes) == 0 {
		return "", errors.New("There is no tag named `<article>` from: " + p.URL.String())
	}
	articleDoc := nodes[0].FirstChild
	plist := htmldoc.ElementsByTag(articleDoc, "p")
	if articleDoc.FirstChild != nil &&
		articleDoc.FirstChild.Data == "div" { // to fetch the summary block
		// body += fmt.Sprintf("\n > %s  \n", plist[0].FirstChild.Data) // redundant summary
		body += fmt.Sprintf("\n > ")
	}
	for _, v := range plist { // the last item is `推荐阅读：`
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.FirstChild != nil && v.FirstChild.Data == "strong" {
			if d := v.FirstChild.FirstChild.Data; !strings.Contains(d, "↓↓↓") ||
				!strings.Contains(d, "点击浏览") {
				body += fmt.Sprintf("\n** %s **  \n", d)
			}
			if t := v.FirstChild.NextSibling; t != nil && t.Type == html.TextNode {
				body += t.Data
			}
		} else {
			ok := true

			for _, a := range v.Parent.Attr {
				if a.Key == "class" {
					switch a.Val {
					// if it is a info for picture, igonre!
					case "sc-bdVaJa iHZvIS":
						ok = false
					// if it is a twitter content, ignore!
					case "twitter-tweet":
						ok = false
					}
				}
			}
			if ok {
				body += v.FirstChild.Data + "  \n"
			}
		}
	}
	body = strings.ReplaceAll(body, "strong", "")
	body = strings.ReplaceAll(body, "** 推荐阅读： **", "")
	return body, nil
}
