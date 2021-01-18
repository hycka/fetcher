package fetcher

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hi20160616/fetcher/internal/htmldoc"
)

func TestSetAndSavePost(t *testing.T) {
	// p := PostFactory("https://www.dwnews.com/经济/60203253")
	// p := PostFactory("https://www.rfa.org/mandarin/Xinwen/6-07082020110802.html") // The wrong one
	// p := PostFactory("https://www.zaobao.com/realtime/world/story20200825-1079575")
	p := NewPost("https://www.cna.com.tw/news/aopl/202009290075.aspx")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := p.TreatPost(); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}

func TestTreatPost(t *testing.T) {
	tcs := []string{
		// "https://www.boxun.com/news/gb/taiwan/2020/07/202007091815.shtml",
		// "https://www.dwnews.com/经济/60203253",
		// "https://www.dwnews.com/全球/60203234",
		// "https://www.voachinese.com/a/S-Korea-Says-US-Sees-Importance-Of-N-Korea-Talks-Despite-Tension-20200709/5496028.html",
		// "https://www.rfa.org/mandarin/yataibaodao/shaoshuminzu/gf1-07092020074142.html",
		// "https://www.rfa.org/mandarin/Xinwen/6-07082020110802.html",
		// "https://www.zaobao.com/realtime/world/story20200825-1079575",
		// "https://www.zaobao.com/news/world/story20200825-1079477",
		// "https://www.zaobao.com.sg/realtime/world/story20200901-1081441",
		// "https://news.ltn.com.tw/news/world/breakingnews/3278726",
		"https://www.cna.com.tw/news/aopl/202009290075.aspx",
	}
	for _, tc := range tcs {
		p := NewPost(tc)
		err := p.TreatPost()
		if err != nil {
			log.Println(err)
		}
	}
}
