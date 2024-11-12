package reload

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

func TestGetRss(t *testing.T) {
	testFeeds := map[string]string{
		"Not Related! A Big-Braned Podcast": "https://notrelated.xyz/rss",
		"Arch Linux: Recent news updates":   "https://www.archlinux.org/feeds/news/",
		"Path of Exile News":                "https://www.pathofexile.com/news/rss",
		"Starsector":                        "https://fractalsoftworks.com/feed/",
		"ShortFatOtaku on Odysee":           "https://odysee.com/$/rss/@ShortFatOtaku:1",
		"CaravanPalace":                     "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
	}
	cacheTime := 60 * time.Minute
	cachePath := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		t.Fatal("cannot create cache directory")
	}
	header := http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
	reloader := new(StandardReloader)
	for title, url := range testFeeds {
		fmt.Println("first iteration...")
		reloader.GetRss(url, header, cacheTime, cachePath)
		fmt.Println("now try to get contents from cache...")
		feed, items, fromCache, err := reloader.GetRss(url, header, cacheTime, cachePath)
		if err != nil {
			t.Fatal(err)
		}
		if !fromCache {
			t.Fatal("did not read from cache for", url)
		}
		if title != feed.Title {
			t.Errorf("expected: %s, parsed: %s\n", title, feed.Title)
		}
		if len(items) == 0 {
			t.Fatal("no items in feed")
		}
	}
}
