package aljazeera

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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
		return errors.New("p.DOC is nil")
	}
	doc := htmldoc.ElementsByTagAndType(p.DOC, "script", "application/ld+json")
	if doc == nil {
		return errors.New("[-] rfa SetDate err, cannot get target nodes.")
	}
	d := doc[0].FirstChild
	if d.Type != html.TextNode {
		return errors.New("[-] rfa SetDate err, target node have no text.")
	}
	raw := d.Data
	re := regexp.MustCompile(`"date\w*?":\s*?"(.*?)"`)
	rs := re.FindAllStringSubmatch(raw, -1)
	p.Date = rs[1][1] // dateModified -> rs[1][1], datePublished -> rs[0][1]
	//UTC add 8H
	if t, err := add8Hour(p.Date); err == nil {
		p.Date = t
	}
	return nil
}

//UTC + 8H
func add8Hour(u string) (string, error) {
	t, err := time.Parse(time.RFC3339, u)
	if err != nil {
		return "", err
	}
	h, _ := time.ParseDuration("+1h")
	h1 := t.Add(8 * h)
	return h1.Format(time.RFC3339), nil
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
		return fmt.Errorf("err at 69L, there is no element <title>")
	}
	title := n[0].FirstChild.Data
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
	b, err := aljazeera(p)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	u, err := url.QueryUnescape(p.URL.String())
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = h1 + "\n\n" + b + "\n\n原地址：" + u
	return nil
}

func aljazeera(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "wysiwyg wysiwyg--all-content")
	if len(nodes) == 0 {
		return "", errors.New("err at 108L, ElementsByTag match nothing from: " + p.URL.String())
	}
	articleDoc := nodes[0]
	plist := htmldoc.ElementsByTag(articleDoc, "h3", "p")

	for _, v := range plist {
		if v.FirstChild != nil {
			if v.Data == "h3" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.Data)
			} else if v.FirstChild.Data == "b" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.FirstChild.Data)
			} else {
				body += v.FirstChild.Data + "  \n"
			}
		}
	}
	return body, nil
}
