package nvimboat

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/mangapill"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/pixiv"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

// Custom reloaders can be defined here.
// Use a regex as the key value to decide, when your reloader should be used.
var CustomReload = map[string]reload.Reloader{
	`https://mangapill\.com`:            new(mangapill.MangapillReloader),
	`https://(?:www\.)?www\.pixiv\.net`: new(pixiv.PixivReloader),
}

func ReloadFeeds(feedUrls []string) (err error) {
	standardReloader := new(reload.StandardReloader)
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
		return
	}
	defer dbh.DB.Close()
	var useCustomReloader bool
	urlReloaderMap := make(map[string]reload.Reloader)

findReloader:
	for _, feedUrl := range feedUrls {
	matchCustomReloader:
		for urlPattern, reloader := range CustomReload {
			useCustomReloader, err = regexp.MatchString(urlPattern, feedUrl)
			if err != nil {
				return errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
			}
			if !useCustomReloader {
				continue matchCustomReloader
			}
			urlReloaderMap[feedUrl] = reloader
			continue findReloader
		}
		urlReloaderMap[feedUrl] = standardReloader
	}

	for url, reloader := range urlReloaderMap {
		err = reloadFeed(url, dbh, reloader)
		if err != nil {
			Log(err)
		}
	}
	Log("Finished reloading.")
	return
}

func reloadFeed(url string, dbh rssdb.DbHandle, reloader reload.Reloader) (err error) {
	var (
		rss_feed       *rssdb.InsertFeedParams
		rss_items      map[string]*rssdb.InsertArticleParams
		fromCache      bool
	)
	header := http.Header{
		"User-Agent": {UserAgent},
	}
	rss_feed, rss_items, fromCache, err = reloader.GetRss(url, header, CacheTime, CachePath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
		return
	}
	if fromCache {
		Log("Loaded " + url + " from cache")
	} else {
		Log("Requested " + url)
	}
	err = reloader.UpdateFeed(dbh, *rss_feed, rss_items)
	if err != nil {
		Log(err)
	}
	return
}
