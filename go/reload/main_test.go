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
	"https://einyakamnil.xyz",
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
	header := http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
	for _, url := range urls {
		fmt.Println("first iteration...")
		getUrlContent(url, header, cacheTime, cachePath)
		fmt.Println("now try to get contents from cache...")
		_, fromCache, err := getUrlContent(url, header, cacheTime, cachePath)
		if err != nil {
			t.Fatal(err)
		}
		if !fromCache {
			t.Fatalf("did not read from cache for: %s\n", url)
		}
	}
}
