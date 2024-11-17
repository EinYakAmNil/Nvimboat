package mangapill

import (
	"fmt"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type MangapillReloader struct {
}

func (mr *MangapillReloader) UpdateFeed(
	dbh rssdb.DbHandle,
	url string,
	cacheTime time.Duration,
	cacheDir string,
	addFeed bool,
) (newFeed rssdb.RssFeed, err error) {
	fmt.Println("Using MangapillReloader for:", url)
	return
}

func (mr *MangapillReloader) GetRss(url string,
	cacheTime time.Duration,
	cacheDir string,
) (feed *rssdb.RssFeed, items []*rssdb.RssItem, fromCache bool, err error) {
	return
}
