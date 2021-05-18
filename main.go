package main

import (
	"log"
	"time"

	"github.com/hi20160616/fetcher/internal/fetcher"
)

func main() {
	worklist := []string{
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
		"https://www.bbc.com/zhongwen/simp/topics/ck2l9z0em07t",
		"https://chinese.aljazeera.net/news",
		"https://cn.reuters.com/assets/jsonWireNews?limit=100",
		"http://cn.kabar.kg/news/",
		"https://ucpnz.co.nz/",
		"https://www.dw.com/zh/%E5%9C%A8%E7%BA%BF%E6%8A%A5%E5%AF%BC/s-9058",
	}
	for {
		fetcher.BreadthFirst(fetcher.Crawl, worklist)
		log.Println("Sleep 10 minutes...")
		time.Sleep(7 * time.Minute)
	}
}
