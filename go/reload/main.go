package reload

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/mmcdole/gofeed"
)

type Reloader interface {
	UpdateFeed(
		url string,
		header http.Header,
		cacheTime time.Duration,
		cacheDir string,
		dbPath string,
	) (err error)

	GetRss(url string,
		header http.Header,
		cacheTime time.Duration,
		cacheDir string,
	) (feed *rssdb.RssFeed, items map[string]*rssdb.RssItem, fromCache bool, err error)

	GetFeed(feedurl string) (feed *rssdb.RssFeed, err error)

	ListFeeds(condition string) (feeds []*rssdb.RssFeed, err error)

	ListArticles(feedurl string) (articles map[string]*rssdb.RssItem, err error)

	AddFeed(
		feed rssdb.CreateFeedParams,
		queries *rssdb.Queries,
		ctx context.Context,
	) (err error)

	AddArticles(
		articles map[string]*rssdb.AddArticlesParams,
		queries *rssdb.Queries,
		ctx context.Context,
	) (err error)
}

type StandardReloader struct{}

func (sr *StandardReloader) UpdateFeed(
	url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
	dbPath string,
) (err error) {
	return
}

// Requests the URL if not found in cacheDir or if the modification time of the cache file is too old.
// The request will be cached in cacheDir.
// Indicates with the return value fromCache if cache was used.
func (sr *StandardReloader) GetRss(url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
) (
	feed *rssdb.RssFeed, items map[string]*rssdb.RssItem, fromCache bool, err error,
) {
	var content []byte
	var reqErr error
	rssParser := gofeed.NewParser()
	cachePath := path.Join(cacheDir, hashUrl(url))
	fileStats, err := os.Stat(cachePath)
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

func (sr *StandardReloader) ListFeeds(condition string) (feeds []*rssdb.RssFeed, err error) {
	return
}

func (sr *StandardReloader) ListArticles(feedurl string) (articles map[string]*rssdb.RssItem, err error) {
	return
}

func (sr *StandardReloader) AddFeed(
	feed rssdb.CreateFeedParams,
	queries *rssdb.Queries,
	ctx context.Context,
) (err error) {
	return
}

func (sr *StandardReloader) AddArticles(
	articles map[string]*rssdb.AddArticlesParams,
	queries *rssdb.Queries,
	ctx context.Context,
) (err error) {
	return
}

func requestUrl(url string, header http.Header) (content []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		err = fmt.Errorf("requestUrl: failed to request %s, status code: %d\n", url, resp.StatusCode)
		return
	}
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	return
}

func cacheUrl(url string, cacheDir string, content []byte) (err error) {
	fileName := hashUrl(url)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		err = fmt.Errorf("cacheUrl: %w", err)
		return
	}
	err = os.WriteFile(path.Join(cacheDir, fileName), content, 0644)
	if err != nil {
		err = fmt.Errorf("cacheUrl: %w", err)
		return
	}
	return
}

func hashUrl(url string) (fileName string) {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	hashBytes := hasher.Sum(nil)
	fileName = hex.EncodeToString(hashBytes)
	fileName = path.Clean(fileName)
	return
}
