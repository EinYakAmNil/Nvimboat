package reload

import (
	"context"
	"database/sql"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Reloader interface {
		UpdateFeed(
			dbh DbHandle,
			url string,
			cacheTime time.Duration,
			cacheDir string,
			addFeed bool,
		) (newFeed rssdb.RssFeed, err error)

		GetRss(url string,
			cacheTime time.Duration,
			cacheDir string,
		) (feed *rssdb.RssFeed, items []*rssdb.RssItem, fromCache bool, err error)
	}

	DbHandle struct {
		DB      *sql.DB
		Ctx     context.Context
		Queries *rssdb.Queries
	}
)
