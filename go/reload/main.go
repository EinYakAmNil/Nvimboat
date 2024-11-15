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
		UpdateFeed(url string, cacheTime time.Duration, cacheDir string, dbh DbHandle) (err error)

		GetRss(url string,
			cacheTime time.Duration,
			cacheDir string,
		) (feed *rssdb.RssFeed, items []*rssdb.RssItem, fromCache bool, err error)

		GetFeed(feedurl string) (feed *rssdb.RssFeed, err error)

		ListFeeds(condition string) (feeds []*rssdb.RssFeed, err error)

		ListArticles(feedurl string) (articles []*rssdb.RssItem, err error)

		AddFeed(feed rssdb.CreateFeedParams, dbh DbHandle) (err error)

		AddArticles(articles []*rssdb.AddArticleParams, dbh DbHandle) (err error)
	}

	DbHandle struct {
		DB      *sql.DB
		Ctx     context.Context
		Queries *rssdb.Queries
	}
)
