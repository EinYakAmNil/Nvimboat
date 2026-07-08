package reload

import (
	"database/sql"
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
	feed rssdb.InsertFeedParams,
	items map[string]*rssdb.InsertArticleParams,
) (err error) {
	_, err = dbh.Queries.GetFeed(dbh.Ctx, feed.Rssurl)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = insertFeed(feed, dbh)
	}
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpdateFeed"))
		return
	}

	knownGuids, err := dbh.Queries.GetFeedGuids(dbh.Ctx, feed.Rssurl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpdateFeed"))
		return
	}
	guidMap := make(map[string]any, len(knownGuids))
	for _, g := range knownGuids {
		guidMap[g] = struct{}{}
	}
	newArticles := make(map[string]*rssdb.InsertArticleParams)
	for g, item := range items {
		if _, ok := guidMap[g]; !ok {
			newArticles[g] = item
			newArticles[g].Unread = 1
		}
	}

	tx, err := dbh.DB.BeginTx(dbh.Ctx, nil)
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpdateFeed"))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	qtx := dbh.Queries.WithTx(tx)
	for _, item := range newArticles {
		err = qtx.InsertArticle(dbh.Ctx, *item)
		if err != nil {
			err = errors.Join(err, errors.New("reload/StandardReloader.UpdateFeed"))
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpdateFeed"))
		return
	}
	return
}

// Requests the URL if not found in cacheDir or if the modification time of the cache file is too old.
// The request will be cached in cacheDir.
// Indicates with the return value fromCache if cache was used.
func (sr *StandardReloader) GetRss(
	url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
) (
	feed *rssdb.InsertFeedParams,
	items map[string]*rssdb.InsertArticleParams,
	fromCache bool,
	err error,
) {
	var (
		content []byte
		reqErr  error
	)
	cachePath := path.Join(cacheDir, hashUrl(url))
	fileStats, err := os.Stat(cachePath)
	// Check if file exists and if the modification time is within cache duration.
	if err != nil || time.Since(fileStats.ModTime()) > cacheTime {
		fromCache = false
		err = fmt.Errorf("%s is not cached", url)
		content, reqErr = requestUrl(url, header)
		if reqErr != nil {
			reqErr = errors.Join(err, reqErr)
			reqErr = errors.Join(reqErr, errors.New("reload/StandardReloader.GetRss"))
			return
		}
		log.Println("requested", url)
		err = cacheUrl(url, cacheDir, content)
		if err != nil {
			err = errors.Join(err, errors.New("reload/StandardReloader.GetRss"))
			return
		}
		log.Println("cached", url)
	} else {
		fromCache = true
		log.Printf("reading %s from cache\n", url)
		content, err = os.ReadFile(cachePath)
		if err != nil {
			err = errors.Join(err, errors.New("reload/StandardReloader.GetRss"))
			return
		}
	}
	feed, items, err = parser.ParseFeed(content, url)
	if err != nil {
		err = fmt.Errorf("parser.ParseFeed: %w\n"+
			"reload/StandardReloader.GetRss", err,
		)
		return
	}
	return
}
