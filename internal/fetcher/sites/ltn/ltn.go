package ltn

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/liuzl/gocc"
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
	p.Err = transform(p)
	return p.Err
}

func setDate(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	metas := htmldoc.MetasByProperty(p.DOC, "article:published_time")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("SetData got nothing.")
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
	if strings.Contains(title, "- 娛樂") ||
		strings.Contains(title, "- 食譜") ||
		strings.Contains(title, "- 地產") ||
		strings.Contains(title, "- 體育") ||
		strings.Contains(title, "- 地方") ||
		strings.Contains(title, "- 蒐奇") ||
		strings.Contains(title, "- 社會") ||
		strings.Contains(title, "- 生活") ||
		strings.Contains(title, "- 时尚") ||
		strings.Contains(title, "- 健康") ||
		strings.Contains(title, "- 汽車") ||
		strings.Contains(title, "- 財經") {
		return errors.New("ignore post on purpose: " + p.URL.String())
	}
	title = strings.ReplaceAll(title, " - 自由時報電子報", "")
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
	b, err := ltn(p)
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

func ltn(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.Raw == nil {
		return "", fmt.Errorf("p.Raw is nil")
	}
	raw := p.Raw
	// Fetch content nodes
	r := htmldoc.DivWithAttr2(raw, "data-desc", "內容頁")
	ps := [][]byte{}
	b := bytes.Buffer{}
	re := regexp.MustCompile(`<p>(.*?)</p>`)
	for _, v := range re.FindAllSubmatch(r, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", fmt.Errorf("no <p> matched")
	}
	for _, p := range ps {
		b.Write(p)
		b.Write([]byte("  \n"))
	}
	body := b.String()
	re = regexp.MustCompile(`「`)
	body = re.ReplaceAllString(body, "“")
	re = regexp.MustCompile(`」`)
	body = re.ReplaceAllString(body, "”")
	re = regexp.MustCompile(`<a.*?>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`</a>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<script.*?</script>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<blockquote.*?</blockquote>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<iframe.*?</iframe>`)
	body = re.ReplaceAllString(body, "")

	return body, nil
}

//transform HANZI
func transform(p *Post) error {
	tw2s, err := gocc.New("tw2s")
	if err != nil {
		p.Err = err
		return err
	}
	//transform title
	in := p.Title
	out, err := tw2s.Convert(in)
	if err != nil {
		p.Err = err
		return err
	}
	p.Title = out
	//transform body
	in = p.Body
	out, err = tw2s.Convert(in)
	if err != nil {
		p.Err = err
		return err
	}
	p.Body = out
	return p.Err
}
