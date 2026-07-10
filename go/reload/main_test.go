package reload

import (
	"errors"
	"fmt"
	"net/http"
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
		// "Starsector":                        "https://fractalsoftworks.com/feed/",
		"ShortFatOtaku on Odysee": "https://odysee.com/$/rss/@ShortFatOtaku:1",
		"CaravanPalace":           "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
		"依云's Blog":               "https://blog.lilydjwg.me/feed",
	}
	cacheTime = 60 * time.Minute
	cacheDir  = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	dbPath    = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "tag-test.db")
	header    = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
)

func TestGetRss(t *testing.T) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatal("cannot create cache directory")
	}
	for title, url := range testFeeds {
		fmt.Println("first iteration...")
		GetRss(url, header, cacheTime, cacheDir)
		fmt.Println("now try to get contents from cache...")
		feed, items, fromCache, err := GetRss(url, header, cacheTime, cacheDir)
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
	dbh, err := rssdb.ConnectDb(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer dbh.DB.Close()
	for _, url := range testFeeds {
		feed, items, _, err := GetRss(url, header, cacheTime, cacheDir)
		if err != nil {
			err = errors.Join(err, errors.New("reload/TestUpdateFeeds"))
			t.Fatal(err)
		}
		err = UpdateFeed(dbh, *feed, items)
		if err != nil {
			err = errors.Join(err, errors.New("reload/TestUpdateFeeds"))
			t.Fatal(err)
		}
	}
}

func TestGetFeed(t *testing.T) {
	dbh, err := rssdb.ConnectDb(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer dbh.DB.Close()
	for title, url := range testFeeds {
		feed, err := dbh.Queries.GetFeed(dbh.Ctx, url)
		if err != nil {
			fmt.Println("Error querying", title, "with", url)
			fmt.Println("Query returned: ", feed)
			t.Fatal(err)
		}
	}
}
