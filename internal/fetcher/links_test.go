package fetcher

import (
	"net/url"
	"testing"
)

func TestKickOutLinksMatchPath(t *testing.T) {
	// link := "https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"
	beKick := []string{"https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"}
	// path := "/" + url.QueryEscape("视觉") + "/"
	path := url.QueryEscape("视觉")
	KickOutLinksMatchPath(&beKick, path)
	if len(beKick) != 0 {
		t.Errorf("want: len(beKick) == 0, got: len(beKick) == %d", len(beKick))
	}
}

func TestSetLinks(t *testing.T) {
	// u, err := url.Parse("https://news.ltn.com.tw/list/breakingnews")
	// assertLinks := []string{
	//         "https://news.ltn.com.tw/news/society/breakingnews/3278253",
	//         "https://news.ltn.com.tw/news/society/breakingnews/3278250",
	//         "https://news.ltn.com.tw/news/politics/breakingnews/3278225",
	//         "https://news.ltn.com.tw/news/politics/breakingnews/3278170",
	// }
	// u, err := url.Parse("https://www.cna.com.tw/list/aall.aspx")
	// assertLinks := []string{
	//         "https://www.cna.com.tw/news/aopl/202009290075.aspx",
	//         "https://www.cna.com.tw/news/firstnews/202009290051.aspx",
	//         "https://www.cna.com.tw/news/acn/202009290063.aspx",
	//         "https://www.cna.com.tw/news/aipl/202009290055.aspx",
	// }
	u, err := url.Parse("https://www.bbc.com/zhongwen/simp/topics/ck2l9z0em07t")
	assertLinks := []string{
		"/zhongwen/simp/world-55655858",
		"/zhongwen/simp/world-55653976",
		"/zhongwen/simp/science-55632120",
		"/zhongwen/simp/chinese-news-55635500",
	}
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
	}
	f.SetLinks()
	shot := 0
	for _, link := range f.Links {
		for _, v := range assertLinks {
			if link == v {
				shot++
			}
		}
	}
	if shot != len(assertLinks) {
		t.Errorf("want: %v, got: %v", len(assertLinks), shot)
	}
}
