package nvimboat

import (
	"fmt"
	"regexp"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/mangapill"
)

var CustomReload = map[string]reload.Reloader{
	"https://mangapill.com": new(mangapill.MangapillReloader),
}

func (nb *Nvimboat) ReloadFeeds(feedUrls []string) (err error) {
	standardReloader := new(reload.StandardReloader)
	dbh, err := reload.ConnectDb(nb.DbPath)
	defer dbh.DB.Close()
	if err != nil {
		err = fmt.Errorf("ReloadFeed: %w", err)
		return
	}
	if err != nil {
		err = fmt.Errorf("ReloadFeed: %w", err)
		return
	}
	var reloadErr error
reloadFeed:
	for _, feedUrl := range feedUrls {
		for urlPattern, reloader := range CustomReload {
			ok, err := regexp.MatchString(urlPattern, feedUrl)
			if err != nil {
				err = fmt.Errorf("ReloadFeeds: %w. Feed url: %s, pattern: %s", err, feedUrl, urlPattern)
				return err
			}
			if ok {
				reloadErr = reloader.UpdateFeed(feedUrl, nb.CacheTime, nb.CachePath, dbh)
				if reloadErr != nil {
					nb.Log(reloadErr)
				}
				continue reloadFeed
			}
		}
		reloadErr = standardReloader.UpdateFeed(feedUrl, nb.CacheTime, nb.CachePath, dbh)
		if reloadErr != nil {
			nb.Log(reloadErr)
		}
	}
	return
}
