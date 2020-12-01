package main

import (
	"log"
	"strconv"
	"time"

	"github.com/hi20160616/fetcher/internal/fetcher"
)

func main() {
	year := strconv.Itoa(time.Now().Year())
	sites := []string{
		// expand fetch range need update cna.go function: setTitle
		"https://www.cna.com.tw/list/aopl.aspx",           // 国际
		"https://news.ltn.com.tw/list/breakingnews/world", // 国际
		"https://www.zaobao.com/realtime/world",
		"https://www.zaobao.com/news/world",
		"https://www.zaobao.com/realtime/china",
		"https://www.zaobao.com/news/china",
		"https://www.dwnews.com",
		"https://www.dwnews.com/issue/10062",
		"https://www.dwnews.com/zone/10000117",
		"https://www.dwnews.com/zone/10000118",
		"https://www.dwnews.com/zone/10000119",
		"https://www.dwnews.com/zone/10000120",
		"https://www.dwnews.com/zone/10000123",
		"https://www.voachinese.com",
		"https://www.voachinese.com/z/1739",
		"https://www.rfa.org/mandarin/",
		"https://www.rfa.org/mandarin/Xinwen/story_archive?year=" + year,
		"https://www.rfa.org/mandarin/yataibaodao/story_archive?year=" + year,
		"https://www.boxun.com/rolling.shtml",
	}
	for {
		fetcher.BreadthFirst(fetcher.Crawl, sites)
		log.Println("Hold a sec ...")
		time.Sleep(5 * time.Minute)
	}
}
