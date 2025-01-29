package reload

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/parser"
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type StandardReloader struct{}

func (sr *StandardReloader) UpdateFeed(
	dbh rssdb.DbHandle,
	feedurl string,
	cacheTime time.Duration,
	cachePath string,
	addFeed bool,
) (newFeed rssdb.RssFeed, err error) {
	knownArticles, err := dbh.Queries.AllArticles(dbh.Ctx)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	feed, items, _, err := sr.GetRss(feedurl, cacheTime, cachePath)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	if addFeed {
		feedParams := rssdb.CreateFeedParams{
			Rssurl: feedurl,
			Title:  feed.Title,
			Url:    feed.Url,
		}
		newFeed, err = sr.AddFeed(feedParams, dbh)
		if err != nil {
			err = fmt.Errorf("UpdateFeed: %w", err)
			return
		}
	}
	var (
		itemsParams = []*rssdb.AddArticleParams{}
	)
	for _, i := range items {
		if _, ok := knownArticles[i.Guid]; !ok {
			itemsParams = append(itemsParams, &rssdb.AddArticleParams{
				Guid:            i.Guid,
				Title:           i.Title,
				Author:          i.Author,
				Url:             i.Url,
				Feedurl:         feedurl,
				Pubdate:         i.Pubdate,
				Content:         i.Content,
				Unread:          i.Unread,
				EnclosureUrl:    i.EnclosureUrl,
				Flags:           i.Flags,
				ContentMimeType: i.ContentMimeType,
			})
		}
	}
	err = sr.AddArticles(itemsParams, feedurl, dbh)
	if err != nil {
		err = fmt.Errorf("UpdateFeed: %w", err)
		return
	}
	return
}

// Requests the URL if not found in cacheDir or if the modification time of the cache file is too old.
// The request will be cached in cacheDir.
// Indicates with the return value fromCache if cache was used.
func (sr *StandardReloader) GetRss(url string,
	cacheTime time.Duration,
	cacheDir string,
) (
	feed *rssdb.RssFeed, items map[string]*rssdb.RssItem, fromCache bool, err error,
) {
	var (
		content []byte
		reqErr  error
	)
	header := http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"},
	}
	rssParser := parser.ParseYtFeed
	cachePath := path.Join(cacheDir, hashUrl(url))
	fileStats, err := os.Stat(cachePath)
	// Check if file exists and if the modification time is within cache duration.
	if err != nil || time.Now().Sub(fileStats.ModTime()) > cacheTime {
		err = fmt.Errorf("%s is not cached", url)
		content, reqErr = requestUrl(url, header)
		if reqErr != nil {
			reqErr = fmt.Errorf("GetRss: %w", errors.Join(err, reqErr))
			return nil, nil, false, reqErr
		}
		log.Println("requested", url)
		err = cacheUrl(url, cacheDir, content)
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, nil, false, err
		}
		log.Println("cached", url)
	} else {
		log.Printf("reading %s from cache\n", url)
		content, err = os.ReadFile(cachePath)
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, nil, false, err
		}
	}
	feedParsed, err := rssParser.ParseString(string(content))
	feed = &rssdb.RssFeed{
		Rssurl: feedParsed.FeedLink,
		Url:    feedParsed.Link,
		Title:  feedParsed.Title,
	}
	if err != nil {
		err = fmt.Errorf("GetRss: %w", err)
		return nil, nil, true, err
	}
	items = make(map[string]*rssdb.RssItem)
	var author string
	for _, item := range feedParsed.Items {
		if len(item.Authors) > 0 {
			author = item.Authors[0].Name
		} else {
			author = ""
		}
		items[item.Link] = &rssdb.RssItem{
			Guid:    item.GUID,
			Title:   item.Title,
			Author:  author,
			Url:     item.Link,
			Feedurl: feedParsed.FeedLink,
			Pubdate: item.PublishedParsed.Unix(),
			Content: item.Content,
			Unread:  1,
		}
	}
	return feed, items, true, err
}

func (sr *StandardReloader) AddFeed(feed rssdb.CreateFeedParams, dbh rssdb.DbHandle) (newFeed rssdb.RssFeed, err error) {
	newFeed, err = dbh.Queries.CreateFeed(dbh.Ctx, feed)
	if err != nil {
		err = fmt.Errorf("AddFeed: %w", err)
		return
	}
	return
}

func (sr *StandardReloader) AddArticles(articles []*rssdb.AddArticleParams, feedUrl string, dbh rssdb.DbHandle) (err error) {
	tx, err := dbh.DB.Begin()
	if err != nil {
		err = fmt.Errorf("AddArticles: %w", err)
	}
	defer tx.Rollback()
	qtx := dbh.Queries.WithTx(tx)
	for _, a := range articles {
		err = qtx.AddArticle(dbh.Ctx, *a)
		if err != nil {
			err = fmt.Errorf("AddArticle: %w, %s", err, a.Guid)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("AddArticle: %w", err)
		return
	}
	return
}
