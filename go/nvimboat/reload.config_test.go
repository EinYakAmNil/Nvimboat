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
	}
	cacheTime = 60 * time.Minute
	cachePath = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test")
	dbPath    = path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "reload_test.db")
)

func TestReloadFeeds(t *testing.T) {
	nb := new(Nvimboat)
	nb.CacheTime = cacheTime
	nb.CachePath = cachePath
	nb.DbPath = dbPath
	err := ReloadFeeds(nb, testFeeds)
	if err != nil {
		t.Fatal(err)
	}
}
