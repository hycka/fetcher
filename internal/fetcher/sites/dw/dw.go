package dw

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/pkg/errors"
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
	re := regexp.MustCompile(`articleChangeDateShort: "(\d*?)",`)
	rs := re.FindAllSubmatch(p.Raw, -1)
	// verbose judgements for pass panic of index out of range.
	if rs != nil && rs[0] != nil && len(rs[0]) > 1 && rs[0][1] != nil {
		t, err := time.Parse("20060102", string(rs[0][1]))
		if err != nil {
			return err
		}
		p.Date = t.Format(time.RFC3339)
	}
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
	title = title[:strings.Index(title, "|")]
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
	tt := t.Format("# [02.01][1504H] " + p.Title)
	u, err := url.QueryUnescape(p.URL.String())
	if err != nil {
		return errors.WithMessage(err, "dw: dw: setBody: url unescape err on: "+p.URL.String())
	}
	p.Body = tt + "\n\n" + b + "\n\n原地址：" + u
	return nil
}

func dw(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", errors.New("dw: p.DOC is nil")
	}
	if p.Raw == nil {
		return "", errors.New("dw: p.Raw is nil")
	}
	doc := p.DOC
	body := ""

	// Fetch summary
	re := regexp.MustCompile(`<p class="intro">(.*?)</p>`)
	raw := bytes.ReplaceAll(p.Raw, []byte("\n"), []byte(""))
	rs := re.FindAllSubmatch(raw, -1)
	if rs == nil {
		return "", errors.New("dw: dw: intro match nothing.")
	}
	if intro := string(rs[0][1]); intro != "" {
		body += "> " + intro + "  \n\n" // if intro exist, append to body
	}

	// Fetch content
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "longText")
	if len(nodes) == 0 {
		return "", errors.New("dw: L118: nodes fetch error from: " + p.URL.String())
	}

	if nodes[0].FirstChild.NextSibling.Attr[0].Val == "col1" {
		nodes[0].RemoveChild(nodes[0].FirstChild.NextSibling)
	}

	spanMerge := func(n *html.Node) []*html.Node {
		spans := htmldoc.ElementsByTag(n, "span")
		for _, span := range spans {
			if span.FirstChild != nil {
				body += span.FirstChild.Data
				if span.FirstChild.Data != span.LastChild.Data {
					body += span.LastChild.Data
				}
			}
		}
		return spans
	}

	plist := htmldoc.ElementsByTag(nodes[0], "p", "h2")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		} else {
			switch v.Data {
			case "h2":
				body += "\n ** "
				if ss := spanMerge(v); len(ss) == 0 {
					body += v.FirstChild.Data
				}
				body += " **   \n"
			case "p":
				if ss := spanMerge(v); len(ss) == 0 {
					body += v.FirstChild.Data
				}
				body += "  \n"
			default:
				body += v.FirstChild.Data + "  \n"
			}
		}
	}
	body = strings.ReplaceAll(body, "strong  \n", "")
	body = strings.ReplaceAll(body, "em  \n", "")
	body = strings.ReplaceAll(body, " ** \n **   \n", "")
	return body, nil
}
