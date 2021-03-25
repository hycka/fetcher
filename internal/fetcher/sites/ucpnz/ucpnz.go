package ucpnz

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hi20160616/fetcher/internal/htmldoc"
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
	doc := htmldoc.ElementsByTagAndClass(p.DOC, "span", "td-post-date")
	d := []string{}
	if doc == nil {
		return fmt.Errorf("there is no element <time>")
	}
	//focus on node like "<span class="td-post-date"><time class="entry-date updated td-module-date" datetime="2020-11-05T13:30:02+00:00" >2020-11-05</time></span>"
	if doc[0].LastChild.Attr[1].Val != "" {
		d = append(d, doc[0].LastChild.Attr[1].Val)
	}
	if len(d) <= 0 {
		return fmt.Errorf("SetData got nothing.")
	}
	p.Date = d[0]
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
	b, err := ucpnz(p)
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

func ucpnz(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.Raw == nil {
		return "", fmt.Errorf("p.Raw is nil")
	}
	raw := p.Raw
	//td-post-content tagdiv-type
	r := htmldoc.DivWithAttr2(raw, "class", "td-post-content tagdiv-type")
	ps := [][]byte{}
	b := bytes.Buffer{}
	re := regexp.MustCompile(`<p.*?>(.*?)</p>`)
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
	re = regexp.MustCompile(`<strong.*?</strong>`)
	body = re.ReplaceAllString(body, "")

	return body, nil
}

//transform HANZI
func transform(p *Post) error {
	tw2s, err := gocc.New("hk2s")
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
