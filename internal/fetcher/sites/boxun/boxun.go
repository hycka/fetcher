package boxun

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
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
	if p.URL == nil {
		return fmt.Errorf("[-] p.URL is nil")
	}
	rawdate := filepath.Base(p.URL.String())
	var Y, M, D, hh, mm int
	var err error
	if Y, err = strconv.Atoi(rawdate[:4]); err != nil {
		Y = 2000
		fmt.Println(rawdate[:4], "is not a Year., set it to 2000")
	}
	if M, err = strconv.Atoi(rawdate[4:6]); err != nil {
		M = 0
		fmt.Println(rawdate[4:6], "is not a Month., set it to 0")
	}
	if D, err = strconv.Atoi(rawdate[6:8]); err != nil {
		D = 0
		fmt.Println(rawdate[6:8], "is not a Date., set it to 0")
	}
	if hh, err = strconv.Atoi(rawdate[8:10]); err != nil {
		hh = 0
		fmt.Println(rawdate[8:10], "is not a Hour., set it to 0")
	}
	if mm, err = strconv.Atoi(rawdate[10:12]); err != nil {
		mm = 0
		fmt.Println(rawdate[10:12], "is not a Minute., set it to 0")
	}
	if err != nil {
		fmt.Println("err date fetch from url: ", p.URL.String())
	}
	if Y < 1999 || Y > 2499 {
		Y = 2000
		// fmt.Println(Y, "is not a integer of Year, set it to 2000")
	}
	if M < 0 || M >= 12 {
		M = 12
		// fmt.Println(M, "is not a integer of Month, set it to 12")
	}
	if D < 0 || D >= 31 {
		D = 31
		// fmt.Println(D, "is not a integer of Date, set it to 31")
	}
	if hh < 0 || hh > 23 {
		hh = 23
		// fmt.Println(hh, "is not a integer of Hour, set it to 23")
	}
	if mm < 0 || mm > 59 {
		mm = 59
		// fmt.Println("err date fetch from url: ", url)
		// fmt.Println(mm, "is not a integer of Minute, set it to 59")
	}
	p.Date = fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02dZ", Y, M, D, hh, mm, 0)
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
	if err := gears.ConvertToUtf8(&title, "gbk", "utf8"); err != nil {
		return err
	}
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	b, err := Boxun(p)
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

func Boxun(p *Post) (string, error) {
	if p.DOC == nil {
		return "", fmt.Errorf("[-] p.DOC is nil")
	}
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndClass(doc, "td", "F11")
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<td class=F11>` from: " + p.URL.String())
	}
	blist := htmldoc.ElementsNextByTag(nodes[0], "br")
	for _, b := range blist {
		if b.Type != html.TextNode || b.Data == "" {
			continue
		} else {
			body += strings.ReplaceAll(b.Data, "\u00a0", "")
		}
	}
	gears.ConvertToUtf8(&body, "gbk", "utf8")
	return body, nil
}
