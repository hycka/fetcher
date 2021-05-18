package zaobao

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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
	if p.Raw == nil {
		return errors.New("zaobao: setDate: Raw is nil")
	}
	re := regexp.MustCompile(`"dateCreated":\s"(\d\d\d\d)-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)`)
	rs := re.FindAllSubmatch(p.Raw, -1)[0]
	m, d, y, hh, mm := rs[2], rs[3], rs[1], rs[4], rs[5]
	p.Date = fmt.Sprintf("%s-%s-%sT%s:%s:00+08:00", y, m, d, hh, mm)
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
	title = strings.ReplaceAll(title, " | 联合早报网", "")
	title = strings.ReplaceAll(title, " | 早报", "")
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
	b, err := zaobao(p)
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

func zaobao(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "col-lg-12 col-12 article-container")
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
