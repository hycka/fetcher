package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hi20160616/fetcher/internal/fetcher/sites/bbc"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/cna"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/dwnews"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/ltn"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/voachinese"
	"github.com/hi20160616/fetcher/internal/fetcher/sites/zaobao"
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

type Paragraph struct {
	Type    string
	Content string
}

func NewPost(rawurl string) *Post {
	url, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Post{
		Domain: url.Hostname(),
		URL:    url,
	}
}

// TODO: use func init
// PostInit open url and get raw and doc
func (p *Post) PostInit() error {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw, p.DOC = raw, doc
	return nil
}

// TreatPost get post things and set to `p` then save it.
func (p *Post) TreatPost() error {
	// Init post
	if err := p.PostInit(); err != nil {
		return err
	}
	// Set post
	switch p.Domain {
	case "www.dwnews.com":
		post := dwnews.Post(*p)
		if err := dwnews.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.voachinese.com":
		post := voachinese.Post(*p)
		if err := voachinese.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.zaobao.com":
		post := zaobao.Post(*p)
		if err := zaobao.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "news.ltn.com.tw":
		post := ltn.Post(*p)
		if err := ltn.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.cna.com.tw":
		post := cna.Post(*p)
		if err := cna.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.bbc.com":
		post := bbc.Post(*p)
		if err := bbc.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	default:
		return fmt.Errorf("switch no case on: %s", p.Domain)
	}
	// Save post to file
	if err := p.setFilename(); err != nil {
		return err
	}
	if err := p.savePost(); err != nil {
		return err
	}
	return nil
}

func (p *Post) savePost() error {
	folderPath := filepath.Join("wwwroot", p.Domain)
	gears.MakeDirAll(folderPath)
	if p.Filename == "" {
		return errors.New("savePost need a filename, but got none.")
	}
	fpath := filepath.Join(folderPath, p.Filename)
	// !+ rm files with same title
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), p.Title) {
			err = os.Remove(filepath.Join(folderPath, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	// !- rm files with same title
	if p.Body == "" {
		p.Body = "savePost p.Body = \"\""
	}
	err = ioutil.WriteFile(fpath, []byte(p.Body), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) setFilename() error {
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	p.Filename = fmt.Sprintf("[%02d.%02d][%02d%02dH]%s.txt", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	return nil
}
