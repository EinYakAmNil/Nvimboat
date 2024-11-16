package nvimboat

import (
	"fmt"
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

func ReloadFeeds(nb *Nvimboat, feedUrls []string) (err error) {
	standardReloader := new(reload.StandardReloader)
	dbh, err := reload.ConnectDb(nb.DbPath)
	defer dbh.DB.Close()
	if err != nil {
		err = fmt.Errorf("ReloadFeed: %w", err)
		return
	}
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
				err = fmt.Errorf("ReloadFeeds: %w. Feed url: %s, pattern: %s", err, feedUrl, urlPattern)
				return err
			}
			if ok {
				newFeed, reloadErr = reloader.UpdateFeed(dbh, feedUrl, nb.CacheTime, nb.CachePath, addFeed)
				if reloadErr != nil {
					nb.Log(reloadErr)
				}
				if addFeed {
					nb.Log("Added feed:", newFeed.Url)
				}
				addFeed = false
				continue reloadFeed
			}
		}
		newFeed, reloadErr = standardReloader.UpdateFeed(dbh, feedUrl, nb.CacheTime, nb.CachePath, addFeed)
		if reloadErr != nil {
			nb.Log(reloadErr)
		}
		if addFeed {
			nb.Log("Added feed:", newFeed.Url)
		}
		addFeed = false
	}
	return
}
