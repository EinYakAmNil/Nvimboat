package reload

import (
	"crypto/sha1"
	"errors"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func requestUrl(url string, header http.Header) (content []byte, err error) {
	client := new(http.Client)
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

func insertFeed(feed rssdb.InsertFeedParams, dbh rssdb.DbHandle) (newFeed rssdb.RssFeed, err error) {
	tx, err := dbh.DB.BeginTx(dbh.Ctx, nil)
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpsertFeed"))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	qtx := dbh.Queries.WithTx(tx)
	newFeed, err = qtx.InsertFeed(dbh.Ctx, feed)
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpsertFeed"))
		return
	}
	err = tx.Commit()
	if err != nil {
		err = errors.Join(err, errors.New("reload/StandardReloader.UpsertFeed"))
		return
	}
	return
}
