package htmldoc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/html"
)

var u, err = url.Parse("https://news.ltn.com.tw/news/world/breakingnews/3277899")

func TestDivWithAttr(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := DivWithAttr(doc, "data-desc", "內容頁")
	plist := ElementsByTag(tc[0], "p")
	for _, v := range plist {
		fmt.Println(v.FirstChild.Data)
	}
}
func TestDivWithAttr2(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := DivWithAttr2(raw, "data-desc", "內容頁")
	fmt.Println(string(tc))
}

func TestElementsByTag(t *testing.T) {
	s, err := ioutil.ReadFile("./bbc.html")
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTag(doc, "main")
	a := ElementsByTag(tc[0], "p")
	for _, v := range a {
		if v.FirstChild != nil {
			fmt.Println(v.FirstChild.Data)
		}
	}
}

func TestElementsByTagAndClass(t *testing.T) {
	s, err := ioutil.ReadFile("./test.html")
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndClass(doc, "div", "paragraph")
	a := ElementsByTag(tc[0], "h2", "p")
	for _, v := range a {
		if v.FirstChild != nil {
			fmt.Println(v.FirstChild.Data)
		}
	}
}

func TestElementsByTagAndClass2(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	raw, _, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndClass2(raw, "div", "wsw")
	fmt.Println(tc)
}

func TestElementsByTagAndId(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndId(doc, "div", "storytext")
	plist := ElementsByTag(tc[0], "p")
	for _, v := range plist {
		if v.FirstChild != nil {
			if v.FirstChild.Data == "b" {
				blist := ElementsByTag(v, "b")
				fmt.Print("**")
				for _, b := range blist {
					fmt.Print(b.FirstChild.Data)
				}
				fmt.Print("**\n")
				// fmt.Println("**" + v.FirstChild.FirstChild.Data + "**")
			} else {
				fmt.Println(v.FirstChild.Data)
			}
		}
	}
}

func TestMetaByProperty(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByProperty(doc, "article:modified_time")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020-08-25T09:42:32+08:00"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestMetaByItemprop(t *testing.T) {
	u, err = url.Parse("https://www.cna.com.tw/news/aopl/202009290075.aspx")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByItemprop(doc, "dateModified")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020/09/29 11:49"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestMetaByName(t *testing.T) {
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByName(doc, "parsely-pub-date")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020-07-09T18:04:00+08:00"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}

func TestElementsRmByTag(t *testing.T) {
	s, err := ioutil.ReadFile("./test.html")
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	ElementsRmByTag(doc, "br")
	n := ElementsByTag(doc, "br")
	if len(n) > 0 {
		t.Errorf("want 0, got: %v", len(n))
		fmt.Println(doc)
	}
}
