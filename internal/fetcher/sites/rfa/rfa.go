package rfa

import (
	"bytes"
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
	p.Date = rs[0][1] // dateModified -> rs[0][1], datePublished -> rs[1][1]
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
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	b, err := Rfa(p)
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

func Rfa(p *Post) (string, error) {
	if p.Raw == nil {
		return "", errors.New("\n[-] p.Raw is nil!\n")
	}
	var ps [][]byte
	var b bytes.Buffer
	var re = regexp.MustCompile(`(?m)<p.*?>(?P<content>.*?)</p>`)
	for _, v := range re.FindAllSubmatch(p.Raw, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		if regexp.MustCompile(`(?m)<video.*?>`).FindAll(p.Raw, -1) != nil {
			return "", errors.New("\n[-] fetcher.FmtBodyRfa() Error: this is a video page.\n")
		}
		return "", errors.New("\n[-] fetcher.FmtBodyRfa() Error: regex matched nothing.\n")
	} else {
		for _, p := range ps {
			b.Write(p)
			b.Write([]byte("  \n"))
		}
	}
	body := string(b.Bytes())
	re = regexp.MustCompile(`(?m)<i.*?</i>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`(?m)<iframe.*?</iframe>`)
	body = re.ReplaceAllString(body, "")
	body = strings.ReplaceAll(body, "\n\n", "\n")
	body = strings.ReplaceAll(body, "<br/>", "")
	body = strings.ReplaceAll(body, "<b>", "** ")
	body = strings.ReplaceAll(body, "</b>", " **  \n")
	body = strings.ReplaceAll(body, "**   **", "")

	return body, nil
}
