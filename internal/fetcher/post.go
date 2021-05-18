package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/aljazeera"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/bbc"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/cna"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/dw"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/dwnews"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/kabar"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/ltn"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/reuters"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/ucpnz"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/voachinese"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/zaobao"
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

type Paragraph struct {
	Type    string
	Content string
}

func NewPost(rawurl string) *Post {
	p := &Post{}
	p.URL, p.Err = url.Parse(rawurl)
	p.Domain = p.URL.Hostname()
	return p
}

// TODO: use func init
// PostInit open url and get raw and doc
func (p *Post) PostInit() error {
	if p.Err != nil {
		return p.Err
	}
	p.Raw, p.DOC, p.Err = htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	return p.Err
}

// RoutePost will switch post to the right dealer.
func (p *Post) RoutePost() error {
	if p.Err != nil {
		return p.Err
	}
	switch p.Domain {
	case "www.dwnews.com":
		post := dwnews.Post(*p)
		p.Err = dwnews.SetPost(&post)
		*p = Post(post)
	case "www.voachinese.com":
		post := voachinese.Post(*p)
		p.Err = voachinese.SetPost(&post)
		*p = Post(post)
	case "www.zaobao.com":
		post := zaobao.Post(*p)
		p.Err = zaobao.SetPost(&post)
		*p = Post(post)
	case "www.zaobao.com.sg":
		post := zaobao.Post(*p)
		p.Err = zaobao.SetPost(&post)
		*p = Post(post)
	case "news.ltn.com.tw":
		post := ltn.Post(*p)
		p.Err = ltn.SetPost(&post)
		*p = Post(post)
	case "www.cna.com.tw":
		post := cna.Post(*p)
		p.Err = cna.SetPost(&post)
		*p = Post(post)
	case "www.bbc.com":
		post := bbc.Post(*p)
		p.Err = bbc.SetPost(&post)
		*p = Post(post)
	case "chinese.aljazeera.net":
		post := aljazeera.Post(*p)
		p.Err = aljazeera.SetPost(&post)
		*p = Post(post)
	case "cn.reuters.com":
		post := reuters.Post(*p)
		p.Err = reuters.SetPost(&post)
		*p = Post(post)
	case "cn.kabar.kg":
		post := kabar.Post(*p)
		p.Err = kabar.SetPost(&post)
		*p = Post(post)
	case "ucpnz.co.nz":
		post := ucpnz.Post(*p)
		p.Err = ucpnz.SetPost(&post)
		*p = Post(post)
	case "www.dw.com":
		post := dw.Post(*p)
		p.Err = dw.SetPost(&post)
		*p = Post(post)
	default:
		return fmt.Errorf("switch no case on: %s", p.Domain)
	}
	return p.Err
}

// TreatPost get post things and set to `p` then save it.
func (p *Post) TreatPost() error {
	// Post prepare
	if p.Err = p.PostInit(); p.Err != nil {
		return p.Err
	}
	if p.Err = p.RoutePost(); p.Err != nil {
		return p.Err
	}

	// Post storage
	if p.Err = p.setFilename(); p.Err != nil {
		return p.Err
	}
	p.Err = p.savePost()

	return p.Err
}

func (p *Post) savePost() error {
	if p.Err != nil {
		return p.Err
	}
	folderPath := filepath.Join("wwwroot", p.Domain)
	gears.MakeDirAll(folderPath)
	if p.Filename == "" {
		return errors.New("savePost need a filename, but got nil.")
	}
	fpath := filepath.Join(folderPath, p.Filename)
	// !+ rm files with same title
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && p.Title != "" && strings.Contains(f.Name(), p.Title) {
			p.Err = os.Remove(filepath.Join(folderPath, f.Name()))
		}
	}
	// !- rm files with same title
	if p.Body == "" {
		p.Body = "savePost p.Body = \"\""
	}
	p.Err = ioutil.WriteFile(fpath, []byte(p.Body), 0644)
	return p.Err
}

func (p *Post) setFilename() error {
	if p.Err != nil {
		return p.Err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	p.Filename = fmt.Sprintf("[%02d.%02d][%02d%02dH]%s.txt", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	return nil
}
