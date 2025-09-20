package nvimboat

import (
	"errors"
	"regexp"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/mangapill"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

// Custom reloaders can be defined here.
// Use a regex as the key value to decide, when your reloader should be used.
var CustomReload = map[string]reload.Reloader{
	"https://mangapill.com": new(mangapill.MangapillReloader),
}

func ReloadFeeds(feedUrls []string) (err error) {
	standardReloader := new(reload.StandardReloader)
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
		return
	}
	defer dbh.DB.Close()
	var (
		newFeed   rssdb.RssFeed
		reloadErr error
		addFeed   bool
	)
	knownFeeds, err := dbh.Queries.MapFeedUrls(dbh.Ctx)
reloadFeed:
	for _, feedUrl := range feedUrls {
		if !knownFeeds[feedUrl] {
			addFeed = true
		}
		for urlPattern, reloader := range CustomReload {
			ok, err := regexp.MatchString(urlPattern, feedUrl)
			if err != nil {
				return errors.Join(err, errors.New("nvimboat/ReloadFeeds"))
			}
			if ok {
				newFeed, reloadErr = reloader.UpdateFeed(dbh, feedUrl, CacheTime, CachePath, addFeed)
				if reloadErr != nil {
					Log(reloadErr)
				}
				if addFeed {
					Log("Added feed:", newFeed.Url)
				}
				addFeed = false
				continue reloadFeed
			}
		}
		newFeed, reloadErr = standardReloader.UpdateFeed(dbh, feedUrl, CacheTime, CachePath, addFeed)
		if reloadErr != nil {
			Log(reloadErr)
		}
		if addFeed {
			Log("Added feed:", newFeed.Url)
		}
		addFeed = false
	}
	return
}
