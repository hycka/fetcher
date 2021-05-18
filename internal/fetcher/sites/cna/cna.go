package cna

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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
	metas := htmldoc.MetasByItemprop(p.DOC, "dateModified")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("setData got nothing.")
	}
	tY := cs[0][:4]
	tM := cs[0][5:7]
	tD := cs[0][8:10]
	tH := cs[0][11:13]
	tm := cs[0][14:16]
	yy, err := strconv.Atoi(tY)
	mm, err := strconv.Atoi(tM)
	dd, err := strconv.Atoi(tD)
	h, err := strconv.Atoi(tH)
	m, err := strconv.Atoi(tm)
	if err != nil {
		return err
	}
	// China doesn't have daylight saving. It uses a fixed 8 hour offset from UTC.
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	t := time.Date(yy, time.Month(mm), dd, h, m, 0, 0, beijing)
	p.Date = t.Format(time.RFC3339)

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
	if strings.Contains(title, "| 娛樂 |") ||
		strings.Contains(title, "| 政治 |") ||
		strings.Contains(title, "| 兩岸 |") ||
		strings.Contains(title, "| 運動 |") ||
		strings.Contains(title, "| 文化 |") ||
		strings.Contains(title, "| 地方 |") ||
		strings.Contains(title, "| 社會 |") ||
		strings.Contains(title, "| 生活 |") ||
		strings.Contains(title, "| 科技 |") ||
		strings.Contains(title, "| 證券 |") ||
		strings.Contains(title, "| 產經 |") {
		return errors.New("ignore post on purpose: " + p.URL.String())
	}
	title = strings.ReplaceAll(title, " | 中央社 CNA", "")
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
	b, err := cna(p)
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

func cna(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.DOC == nil {
		return "", fmt.Errorf("p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "div", "paragraph")
	if len(nodes) == 0 {
		return "", errors.New("There is no element class is paragraph` from: " + p.URL.String())
	}
	n := nodes[0]
	plist := htmldoc.ElementsByTag(n, "h2", "p")
	for _, v := range plist {
		if v.FirstChild != nil {
			body += v.FirstChild.Data + "  \n"
		}
	}

	body = strings.ReplaceAll(body, "「", "“")
	body = strings.ReplaceAll(body, "」", "”")
	body = strings.ReplaceAll(body, "</a>", "")

	re := regexp.MustCompile(`<a.*?>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<script.*?</script>`)
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
