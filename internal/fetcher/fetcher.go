package fetcher

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/hi20160616/gears"
)

type Fetcher struct {
	Entrance *url.URL
	Links    []string
	Err      error
}

func NewFetcher(site string) *Fetcher {
	u, err := url.Parse(site)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Fetcher{
		Entrance: u,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func BreadthFirst(f func(item string), worklist []string) {
	for _, item := range worklist {
		f(item)
	}
}

func Crawl(_url string) {
	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			log.Println(e)
			PanicLog(e)
		}
	}()
	f := NewFetcher(_url)
	log.Printf("[*] Deal with: [%s]\n", _url)
	log.Println("[*] Fetch links ...")
	if err := f.GetLinks(); err != nil { // f.Links update to the _url website is.
		log.Println(err)
		// if links cannot fetch sleep 1 minute then continue
		time.Sleep(1 * time.Minute)
		// continue // only useful by goroutine
		return
	}

	// GetNews then compare via md5 and Save or Rewrite news exist
	log.Println("[*] Get news ...")
	for _, link := range f.Links {
		post := NewPost(link)
		if err := post.TreatPost(); err != nil {
			errMsg := "TreatPost error occur on: " + link
			log.Println(errMsg)
			log.Println(err)
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", f.Entrance.Hostname()), 3)
}

// DelRoutine remove files in folder days ago
func DelRoutine(folder string, n int) error {
	if !gears.Exists(folder) {
		fmt.Printf("\nDelRoutine() err: Folder(%s) does not exist.\n", folder)
		return nil
	}
	// append files to d not be removed
	var c, d, filelist []string
	for i := -1; i < n; i++ {
		a := time.Now().AddDate(0, 0, -i)
		b := fmt.Sprintf("[%02d.%02d]", a.Month(), a.Day())
		c, _ = gears.GetPrefixedFiles(folder, b)
		d = append(d, c...)
		// get file list and rm files not have prefix b
	}
	d = gears.StrSliceDeDupl(d)
	// append files all in the folder
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filelist = append(filelist, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	dellist := gears.StrSliceDiff(filelist, d)
	// if alright, remove files not in d
	for _, file := range dellist {
		log.Println("Del file: ", file)
		os.Remove(file)
	}
	return nil
}

func PanicLog(_err error) error {
	filePath := "./PanicLog.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("[" + time.Now().Format(time.RFC3339) + "]--------------------------------------\n")
	write.WriteString(_err.Error() + "\n")
	write.Flush()
	return nil
}
