package reload

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	_ "github.com/mattn/go-sqlite3"
)

// seems ugly not to embed, but it has multiple reasons:
// 1. if embeded via injection, then testing would be difficult
// 2. sqlite_sequence and sqlite_stat1 should not be created by users, but sqlc needs them in schema
var createDbSql = `
CREATE TABLE rss_feed ( 
	rssurl VARCHAR(1024) PRIMARY KEY NOT NULL,
	url VARCHAR(1024) NOT NULL,
	title VARCHAR(1024) NOT NULL ,
	lastmodified INTEGER(11) NOT NULL DEFAULT 0,
	is_rtl INTEGER(1) NOT NULL DEFAULT 0,
	etag VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE TABLE rss_item (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	guid VARCHAR(64) NOT NULL,
	title VARCHAR(1024) NOT NULL,
	author VARCHAR(1024) NOT NULL,
	url VARCHAR(1024) NOT NULL,
	feedurl VARCHAR(1024) NOT NULL,
	pubDate INTEGER NOT NULL,
	content VARCHAR(65535) NOT NULL,
	unread INTEGER(1) NOT NULL ,
	enclosure_url VARCHAR(1024),
	enclosure_type VARCHAR(1024),
	enqueued INTEGER(1) NOT NULL DEFAULT 0,
	flags VARCHAR(52),
	deleted INTEGER(1) NOT NULL DEFAULT 0,
	base VARCHAR(128) NOT NULL DEFAULT '',
	content_mime_type VARCHAR(255) NOT NULL DEFAULT '',
	enclosure_description VARCHAR(1024) NOT NULL DEFAULT '',
	enclosure_description_mime_type VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE TABLE google_replay (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	guid VARCHAR(64) NOT NULL,
	state INTEGER NOT NULL,
	ts INTEGER NOT NULL 
);
CREATE INDEX idx_rssurl ON rss_feed(rssurl);
CREATE INDEX idx_guid ON rss_item(guid);
CREATE INDEX idx_feedurl ON rss_item(feedurl);
CREATE INDEX idx_lastmodified ON rss_feed(lastmodified);
CREATE INDEX idx_deleted ON rss_item(deleted);
CREATE TABLE metadata ( 
	db_schema_version_major INTEGER NOT NULL,
	db_schema_version_minor INTEGER NOT NULL 
);
`

func ConnectDb(dbPath string) (dbh DbHandle, err error) {
	dbh.Ctx = context.Background()
	dbh.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		err = fmt.Errorf("ConnectDb: %w", err)
		return
	}
	// only create tables, if the database does not exist yet
	if _, noDbErr := os.Stat(dbPath); noDbErr != nil {
		if _, err = dbh.DB.ExecContext(dbh.Ctx, createDbSql); err != nil {
			err = fmt.Errorf("ConnectDb: %w", err)
			return
		}
	}
	dbh.Queries = rssdb.New(dbh.DB)
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
