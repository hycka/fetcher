package fetcher

import (
	"log"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	for {
		BreadthFirst(Crawl, []string{
			// "https://www.boxun.com/rolling.shtml",
			// "https://www.dwnews.com",
			// "https://www.zaobao.com/realtime/world",
			// "https://www.zaobao.com/news/world",
			// "https://www.voachinese.com",
			// "https://www.rfa.org/mandarin/",
			// "https://news.ltn.com.tw/list/breakingnews",
			// "https://www.cna.com.tw/list/aall.aspx",
		})

		log.Println("Sleep a sec ...")
		time.Sleep(5 * time.Minute) // only useful by goroutine
	}
}
