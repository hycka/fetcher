// Pacage links provides a link-extraction fuction.
package fetcher

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/hi20160616/fetcher/internal/htmldoc"
	"github.com/hi20160616/gears"
)

// LinksInit fetch all links from entrance of f, save to f
func (f *Fetcher) LinksInit() error {
	links, err := htmldoc.ExtractLinks(f.Entrance.String())
	if err != nil {
		log.Printf(`can't extract links from "%s": %s`, f.Entrance.String(), err)
		return err
	}
	links = gears.StrSliceDeDupl(links)
	hostname := f.Entrance.Hostname()
	switch hostname {
	case "www.dwnews.com":
		f.Links = LinksFilter(links, `.*?/.*?/\d{8}/`)
		KickOutLinksMatchPath(&f.Links, "zone")
		KickOutLinksMatchPath(&f.Links, "/"+url.QueryEscape("视觉")+"/")
	case "www.voachinese.com":
		l1 := LinksFilter(links, `.*?/a/\d*?.html`)
		l2 := LinksFilter(links, `.*?/a/.*-.*.html`)
		f.Links = append(l1, l2...)
		KickOutLinksMatchPath(&f.Links, "voaweishi")
	case "www.zaobao.com":
		newsWorld := LinksFilter(links, `.*?/news/world/.*`)
		newsChina := LinksFilter(links, `.*?/news/china/.*`)
		realtimeWorld := LinksFilter(links, `.*?/realtime/world/.*`)
		realtimeChina := LinksFilter(links, `.*?/realtime/china/.*`)
		f.Links = append(append(append(newsWorld, newsChina...), realtimeWorld...), realtimeChina...)
	case "news.ltn.com.tw":
		f.Links = LinksFilter(links, `https://news.*/news/.*`)
		KickOutLinksMatchPath(&f.Links, "/life/")
		KickOutLinksMatchPath(&f.Links, "/society/")
		KickOutLinksMatchPath(&f.Links, "/novelty/")
		KickOutLinksMatchPath(&f.Links, "/local/")
	case "www.cna.com.tw":
		newsFirst := LinksFilter(links, `.*?/news/firstnews/.*`)
		newsWorld := LinksFilter(links, `.*?/news/aopl/.*`)
		// TODO: ignore aipl and acn but this still fetch links?
		newsPolitical := LinksFilter(links, `.*?/news/aipl/.*`)
		newsTW := LinksFilter(links, `.*?/news/acn/.*`)
		f.Links = append(append(append(newsFirst, newsWorld...), newsPolitical...), newsTW...)
	case "www.bbc.com":
		f.Links = LinksFilter(links, `.*?/zhongwen/simp/.*-\d*`)
		KickOutLinksMatchPath(&f.Links, "institutional")
	case "chinese.aljazeera.net":
		f.Links = LinksFilter(links, `.*?\/[A-Za-z]+\/\d{4}\/\d{1,2}\/\d{1,2}\/.*`)
	}
	return nil
}

// kickOutLinksMatchPath will kick out the links match the path,
func KickOutLinksMatchPath(links *[]string, path string) {
	tmp := []string{}
	// path = "/" + url.QueryEscape(path) + "/"
	// path = url.QueryEscape(path)
	for _, link := range *links {
		if !strings.Contains(link, path) {
			tmp = append(tmp, link)
		}
	}
	*links = tmp
}

// TODO: use point to impletement LinksFilter
// LinksFilter is support for SetLinks method
func LinksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}
