package reload

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

var (
	testFeeds = map[string]string{
		"Not Related! A Big-Braned Podcast": "https://notrelated.xyz/rss",
		"Arch Linux: Recent news updates":   "https://www.archlinux.org/feeds/news/",
		"Path of Exile News":                "https://www.pathofexile.com/news/rss",
		"Starsector":                        "https://fractalsoftworks.com/feed/",
		"ShortFatOtaku on Odysee":           "https://odysee.com/$/rss/@ShortFatOtaku:1",
		"CaravanPalace":                     "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
		"依云's Blog":                         "https://blog.lilydjwg.me/feed",
	}
	cacheTime = 60 * time.Minute
	cacheDir  = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	dbPath    = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "getrss_test.db")
)

func TestGetRss(t *testing.T) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatal("cannot create cache directory")
	}
	reloader := new(StandardReloader)
	for title, url := range testFeeds {
		fmt.Println("first iteration...")
		reloader.GetRss(url, cacheTime, cacheDir)
		fmt.Println("now try to get contents from cache...")
		feed, items, fromCache, err := reloader.GetRss(url, cacheTime, cacheDir)
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
	dbh, err := rssdb.ConnectDb(dbPath)
	defer dbh.DB.Close()
	if err != nil {
		t.Fatal(err)
	}
	knownFeeds, err := dbh.Queries.MapFeedUrls(dbh.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	var addFeed bool
	for _, url := range testFeeds {
		if !knownFeeds[url] {
			fmt.Println("Here:", url)
			addFeed = true
		}
		_, err := reloader.UpdateFeed(dbh, url, cacheTime, cacheDir, addFeed)
		if err != nil {
			t.Fatal(err)
		}
		addFeed = false
	}
	allArticles, err := dbh.Queries.AllArticles(dbh.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, url := range testFeeds {
		_, err := reloader.UpdateFeed(dbh, url, cacheTime, cacheDir, false)
		if err != nil {
			t.Fatal(err)
		}
	}
	noChange, err := dbh.Queries.AllArticles(dbh.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(allArticles) != len(noChange) {
		t.Fatalf("item count should not increase: %d -> %d\n", len(allArticles), len(noChange))
	}
}

func TestGetFeed(t *testing.T) {
	dbh, err := rssdb.ConnectDb(dbPath)
	defer dbh.DB.Close()
	if err != nil {
		t.Fatal(err)
	}
	for title, url := range testFeeds {
		feed, err := dbh.Queries.GetFeed(dbh.Ctx, url)
		if err != nil {
			fmt.Println("Error querying", title, "with", url)
			fmt.Println("Query returned: ", feed)
			t.Fatal(err)
		}
	}
}
