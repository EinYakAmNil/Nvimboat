package nvimboat

import (
	"errors"
	"net/http"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func ReloadFeeds(feedUrls []string) (err error) {
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
		return
	}
	defer dbh.DB.Close()

	for _, url := range feedUrls {
		err = reloadFeed(url, dbh)
		if err != nil {
			Log(err)
		}
	}
	Log("Finished reloading.")
	return
}

func reloadFeed(url string, dbh rssdb.DbHandle) (err error) {
	var (
		rss_feed  *rssdb.InsertFeedParams
		rss_items map[string]*rssdb.InsertArticleParams
		fromCache bool
	)
	header := http.Header{
		"User-Agent": {UserAgent},
	}
	rss_feed, rss_items, fromCache, err = reload.GetRss(url, header, CacheTime, CachePath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
		return
	}
	if fromCache {
		Log("Loaded " + url + " from cache")
	} else {
		Log("Requested " + url)
	}
	err = reload.UpdateFeed(dbh, *rss_feed, rss_items)
	if err != nil {
		Log(err)
	}
	return
}
