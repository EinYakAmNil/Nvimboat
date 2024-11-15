package mangapill

import (
	"fmt"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type MangapillReloader struct {
}

func (mr *MangapillReloader) UpdateFeed(
	url string,
	cacheTime time.Duration,
	cacheDir string,
	dbh reload.DbHandle,
) (err error) {
	sr := new(reload.StandardReloader)
	err = sr.UpdateFeed(url, cacheTime, cacheDir, dbh)
	return
}

func (mr *MangapillReloader) GetRss(url string,
	cacheTime time.Duration,
	cacheDir string,
) (feed *rssdb.RssFeed, items []*rssdb.RssItem, fromCache bool, err error) {
	return
}

func (mr *MangapillReloader) GetFeed(feedurl string) (feed *rssdb.RssFeed, err error) {
	return
}

func (mr *MangapillReloader) ListFeeds(condition string) (feeds []*rssdb.RssFeed, err error) {
	return
}

func (mr *MangapillReloader) ListArticles(feedurl string) (articles []*rssdb.RssItem, err error) {
	return
}

func (mr *MangapillReloader) AddFeed(feed rssdb.CreateFeedParams, dbh reload.DbHandle) (err error) {
	_, err = dbh.Queries.CreateFeed(dbh.Ctx, feed)
	if err != nil {
		err = fmt.Errorf("AddFeed: %w, %+v", err, feed)
		return
	}
	return
}

func (mr *MangapillReloader) AddArticles(articles []*rssdb.AddArticleParams, dbh reload.DbHandle) (err error) {
	for _, a := range articles {
		err = dbh.Queries.AddArticle(dbh.Ctx, *a)
		if err != nil {
			err = fmt.Errorf("AddArticles: %w, %+v", err, a)
			return
		}
	}
	return
}
