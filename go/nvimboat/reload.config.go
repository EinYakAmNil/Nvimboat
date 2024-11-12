package nvimboat

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/mangapill"
)

var CustomReload = map[string]reload.Reloader{
	"https://mangapill.com": new(mangapill.MangapillReloader),
}

func (nb *Nvimboat) ReloadFeeds(feedUrls []string) (err error) {
	standardReloader := new(reload.StandardReloader)
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
				reloadErr = reloader.UpdateFeed(feedUrl,
					nb.Config["http_header"].(http.Header),
					nb.Config["cache_time"].(time.Duration),
					nb.Config["cache"].(string),
					nb.Config["db_path"].(string),
				)
				if reloadErr != nil {
					nb.Log(reloadErr)
				}
				continue reloadFeed
			}
		}
		reloadErr = standardReloader.UpdateFeed(feedUrl,
			nb.Config["http_header"].(http.Header),
			nb.Config["cache_time"].(time.Duration),
			nb.Config["cache"].(string),
			nb.Config["db_path"].(string),
		)
		if reloadErr != nil {
			nb.Log(reloadErr)
		}
	}
	return
}
