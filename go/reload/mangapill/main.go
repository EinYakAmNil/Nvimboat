package mangapill

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type MangapillReloader struct {
}

func (mr *MangapillReloader) UpdateFeed(
	dbh rssdb.DbHandle,
	feed rssdb.InsertFeedParams,
	items map[string]*rssdb.InsertArticleParams,
) (err error) {
	fmt.Println("Using MangapillReloader for:", feed.Rssurl)
	return
}

func (mr *MangapillReloader) GetRss(
	url string,
	heaeder http.Header,
	cacheTime time.Duration,
	cacheDir string,
) (feed *rssdb.InsertFeedParams, items map[string]*rssdb.InsertArticleParams, fromCache bool, err error) {
	return
}
