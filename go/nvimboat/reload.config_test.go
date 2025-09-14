package nvimboat

import (
	"os"
	"path"
	"testing"
	"time"
)

var (
	testFeeds = []string{
		"https://notrelated.xyz/rss",
		"https://www.archlinux.org/feeds/news/",
		"https://www.pathofexile.com/news/rss",
		"https://fractalsoftworks.com/feed/",
		"https://odysee.com/$/rss/@ShortFatOtaku:1",
		"https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
		"https://mangapill.com/manga/1817/houseki-no-kuni",
		"https://blog.lilydjwg.me/feed",
	}
	cacheTime = 60 * time.Minute
	cachePath = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	dbPath    = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "reload_test.db")
)

func TestReloadFeeds(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	CacheTime = cacheTime
	CachePath = cachePath
	DbPath = dbPath
	err := ReloadFeeds(testFeeds)
	if err != nil {
		t.Fatal(err)
	}
}
