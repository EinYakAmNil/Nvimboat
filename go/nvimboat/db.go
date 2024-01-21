package nvimboat

import (
	"database/sql"
	// "errors"
	// "fmt"
	"os"
	"path"
	// "time"

	_ "github.com/mattn/go-sqlite3"
)

func (nb *Nvimboat) QueryArticle(url string) (Article, error) {
	var a Article
	return a, nil
}

func (nb *Nvimboat) QueryFeed(url string) (Feed, error) {
	var f Feed
	return f, nil
}

func (nb *Nvimboat) QueryFeeds() ([]Feed, error) {
	var (
		f     Feed
		feeds []Feed
	)
	rows, err := nb.multiRow(feedsMinimalQuery)
	if err != nil {
		return feeds, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		if err != nil {
			return feeds, err
		}

		feeds = append(feeds, f)
	}

	return feeds, nil
}

func (nb *Nvimboat) QueryFilter(query string, inTags, exTags []string) (Filter, error) {
	var f Filter
	return f, nil
}

func (nb *Nvimboat) QueryTags() (TagsPage, error) {
	var tp TagsPage
	tp.TagFeedCount = make(map[string]int)
	tp.Feeds = nb.ConfigFeeds
	for _, feed := range tp.Feeds {
		for _, tag := range feed["tags"].([]any) {
			tp.TagFeedCount[tag.(string)]++
		}
	}
	return tp, nil
}

func (nb *Nvimboat) QueryTagFeeds(tag string) (TagFeeds, error) {
	var tf TagFeeds
	return tf, nil
}

func (nb *Nvimboat) singleRow(query string, args ...any) *sql.Row {
	row := nb.DB.QueryRow(query, args...)
	return row
}

func (nb *Nvimboat) multiRow(query string, args ...any) (*sql.Rows, error) {
	rows, err := nb.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
