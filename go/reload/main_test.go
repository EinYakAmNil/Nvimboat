package reload

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

var urls = []string{
	"https://notrelated.xyz/rss",
	"https://www.archlinux.org/feeds/news/",
	"https://www.pathofexile.com/news/rss",
	"https://fractalsoftworks.com/feed/",
	"https://odysee.com/$/rss/@ShortFatOtaku:1",
	"https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
}

func TestGetUrlContent(t *testing.T) {
	cacheTime := 60 * time.Minute
	cachePath := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	err := os.MkdirAll(cachePath, 755)
	if err != nil {
		t.Fatal("cannot create cache directory")
	}
	header := http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
	for _, url := range urls {
		fmt.Println("first iteration...")
		GetRss(url, header, cacheTime, cachePath)
		fmt.Println("now try to get contents from cache...")
		feed, fromCache, err := GetRss(url, header, cacheTime, cachePath)
		if err != nil {
			t.Fatal(err)
		}
		if !fromCache {
			t.Fatal("did not read from cache for", url)
		}
		fmt.Println(feed.Title)
	}
}
