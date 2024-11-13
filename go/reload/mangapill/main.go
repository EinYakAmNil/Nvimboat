package mangapill

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type MangapillReloader struct {
}

func (mr *MangapillReloader) UpdateFeed(
	url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
	dbPath string,
) (err error) {
	sr := new(reload.StandardReloader)
	err = sr.UpdateFeed(url, header, cacheTime, cacheDir, dbPath)
	return
}

func (mr *MangapillReloader) GetRss(url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
) (feed *rssdb.RssFeed, items map[string]*rssdb.RssItem, fromCache bool, err error) {
	return
}

func (mr *MangapillReloader) GetFeed(feedurl string) (feed *rssdb.RssFeed, err error) {
	return
}

func (mr *MangapillReloader) ListFeeds(condition string) (feeds []*rssdb.RssFeed, err error) {
	return
}

func (mr *MangapillReloader) ListArticles(feedurl string) (articles map[string]*rssdb.RssItem, err error) {
	return
}

func (mr *MangapillReloader) AddFeed(feed rssdb.CreateFeedParams, queries *rssdb.Queries, ctx context.Context) (err error) {
	_, err = queries.CreateFeed(ctx, feed)
	if err != nil {
		err = fmt.Errorf("AddFeed: %w, %+v", err, feed)
		return
	}
	return
}

func (mr *MangapillReloader) AddArticles(
	articles map[string]*rssdb.AddArticlesParams,
	queries *rssdb.Queries,
	ctx context.Context,
) (err error) {
	for _, a := range articles {
		err = queries.AddArticles(ctx, *a)
		if err != nil {
			err = fmt.Errorf("AddArticles: %w, %+v", err, a)
			return
		}
	}
	return
}
