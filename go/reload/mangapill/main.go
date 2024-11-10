package mangapill

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/mmcdole/gofeed"
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
	rss, _, err := mr.GetRss(url, header, cacheTime, cacheDir)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	feed := rssdb.CreateFeedParams{
		Rssurl: rss.FeedLink,
		Title:  rss.Title,
		Url:    rss.Link,
	}
	queries, ctx, err := reload.ConnectDb(dbPath)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	err = mr.AddFeed(feed, queries, ctx)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	return
}

func (mr *MangapillReloader) GetRss(url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
) (feed *gofeed.Feed, fromCache bool, err error) {
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

func (mr *MangapillReloader) AddFeed(feed rssdb.CreateFeedParams, queries *rssdb.Queries, ctx context.Context) (err error) {
	_, err = queries.CreateFeed(ctx, feed)
	if err != nil {
		err = fmt.Errorf("AddFeed: %w, %+v", err, feed)
		return
	}
	return
}

func (mr *MangapillReloader) AddArticles(
	articles []rssdb.AddArticlesParams,
	queries *rssdb.Queries,
	ctx context.Context,
) (err error) {
	for _, a := range articles {
		err = queries.AddArticles(ctx, a)
		if err != nil {
			err = fmt.Errorf("AddArticles: %w, %+v", err, a)
			return
		}
	}
	return
}
