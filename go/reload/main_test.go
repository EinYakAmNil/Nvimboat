package reload

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

var (
	TestFeeds = map[string]string{
		"Not Related! A Big-Braned Podcast": "https://notrelated.xyz/rss",
		"Arch Linux: Recent news updates":   "https://www.archlinux.org/feeds/news/",
		"Path of Exile News":                "https://www.pathofexile.com/news/rss",
		"Starsector":                        "https://fractalsoftworks.com/feed/",
		"ShortFatOtaku on Odysee":           "https://odysee.com/$/rss/@ShortFatOtaku:1",
		"CaravanPalace":                     "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
	}
	CacheTime = 60 * time.Minute
	CacheDir  = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	DbPath    = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "cache.db")
	Header    = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
)

func TestGetRss(t *testing.T) {
	err := os.MkdirAll(CacheDir, 0755)
	if err != nil {
		t.Fatal("cannot create cache directory")
	}
	reloader := new(StandardReloader)
	for title, url := range TestFeeds {
		fmt.Println("first iteration...")
		reloader.GetRss(url, CacheTime, CacheDir)
		fmt.Println("now try to get contents from cache...")
		feed, items, fromCache, err := reloader.GetRss(url, CacheTime, CacheDir)
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

func TestUpdateFeeds(t *testing.T) {
	reloader := new(StandardReloader)
	dbh, err := ConnectDb(DbPath)
	defer dbh.DB.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, url := range TestFeeds {
		err := reloader.UpdateFeed(url, CacheTime, CacheDir, dbh)
		if err != nil {
			t.Fatal(err)
		}
	}
	allItems, err := dbh.Queries.AllArticles(dbh.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, url := range TestFeeds {
		err := reloader.UpdateFeed(url, CacheTime, CacheDir, dbh)
		if err != nil {
			t.Fatal(err)
		}
	}
	noChange, err := dbh.Queries.AllArticles(dbh.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(allItems) != len(noChange) {
		t.Fatalf("item count should not increase: %d -> %d\n", len(allItems), len(noChange))
	}
}

func TestGetFeed(t *testing.T) {
	dbh, err := ConnectDb(DbPath)
	defer dbh.DB.Close()
	if err != nil {
		t.Fatal(err)
	}
	feed, err := dbh.Queries.GetFeed(dbh.Ctx, "https://fractalsoftworks.com/feed/")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(feed)
}
