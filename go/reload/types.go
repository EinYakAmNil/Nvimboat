package reload

import (
	"net/http"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	_ "github.com/mattn/go-sqlite3"
)

type Reloader interface {
	UpdateFeed(
		dbh rssdb.DbHandle,
		feed rssdb.InsertFeedParams,
		items map[string]*rssdb.InsertArticleParams,
	) (err error)

	GetRss(
		url string,
		header http.Header,
		cacheTime time.Duration,
		cacheDir string,
	) (
		feed *rssdb.InsertFeedParams,
		items map[string]*rssdb.InsertArticleParams,
		fromCache bool,
		err error,
	)
}
