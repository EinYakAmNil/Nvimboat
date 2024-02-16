package nvimboat

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(dbpath string) (db *sql.DB, err error) {
	if fileExists(dbpath) {
		db, err = sql.Open("sqlite3", dbpath)
		return
	}
	dbDir := path.Dir(dbpath)
	err = os.MkdirAll(dbDir, os.FileMode(0755))
	if err != nil {
		return
	}
	d, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return
	}
	_, err = d.Exec(createDB)
	return
}

func QueryMain(db *sql.DB, configFeeds []map[string]any, configFilters []map[string]any) (*MainMenu, error) {
	var (
		err        error
		tmp_filter Filter
		mainmenu   = &MainMenu{ConfigFeeds: configFeeds, ConfigFilters: configFilters}
	)
	mainmenu.Feeds, err = QueryFeeds(db)
	if err != nil {
		return mainmenu, err
	}
	mainmenu.Filters, err = parseFilters(configFilters)
	if err != nil {
		return mainmenu, err
	}
	for i, f := range mainmenu.Filters {
		tmp_filter, err = QueryFilter(db, configFeeds, f.Query, f.IncludeTags, f.ExcludeTags)
		if err != nil {
			return mainmenu, err
		}
		mainmenu.Filters[i].UnreadCount = tmp_filter.UnreadCount
		mainmenu.Filters[i].ArticleCount = tmp_filter.ArticleCount
		mainmenu.Filters[i].Articles = tmp_filter.Articles
	}
	return mainmenu, err
}

func QueryFeed(db *sql.DB, feedUrl string) (Feed, error) {
	feed := Feed{RssUrl: feedUrl, UnreadCount: 0}
	row := db.QueryRow(feedQuery, feed.RssUrl)
	err := row.Scan(&feed.Title)
	if err != nil {
		return feed, err
	}
	rows, err := db.Query(feedArticlesQuery, feed.RssUrl)
	if err != nil {
		return feed, err
	}
	defer rows.Close()

	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Guid, &a.Title, &a.Author, &a.Url, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
		if err != nil {
			return feed, err
		}
		feed.Articles = append(feed.Articles, a)
	}
	feed.ArticleCount = len(feed.Articles)
	for _, a := range feed.Articles {
		if a.Unread == 1 {
			feed.UnreadCount++
		}
	}
	return feed, nil
}

func QueryFeeds(db *sql.DB) (feeds []*Feed, err error) {
	rows, err := db.Query(feedsMinimalQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		f := new(Feed)
		err = rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		if err != nil {
			return
		}
		feeds = append(feeds, f)
	}
	return
}

func QueryFilter(db *sql.DB, configFeeds []map[string]any, query string, inTags, exTags []string) (Filter, error) {
	var (
		f Filter
	)
	f.Query = query
	f.IncludeTags = inTags
	f.ExcludeTags = exTags
	urls := filterTags(configFeeds, inTags, exTags)
	if len(urls) == 0 {
		return f, nil
	}
	q := articlesFilterQuery(query, len(urls))
	rows, err := db.Query(q, urls...)
	if err != nil {
		return f, err
	}
	defer rows.Close()
	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Guid, &a.Title, &a.Author, &a.Url, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
		if err != nil {
			return f, nil
		}
		f.ArticleCount++
		if a.Unread == 1 {
			f.UnreadCount++
		}
		f.Articles = append(f.Articles, a)
	}
	return f, nil
}

func QueryArticle(db *sql.DB, url string) (Article, error) {
	var a = Article{Url: url}
	row := db.QueryRow(articleQuery, a.Url)
	err := row.Scan(&a.Guid, &a.Title, &a.Author, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
	if err != nil {
		return a, err
	}
	return a, nil
}

func QueryTagFeeds(db *sql.DB, tag string, configFeeds []map[string]any) (TagFeeds, error) {
	var (
		tf       TagFeeds
		feedurls []any
	)
	tf.Tag = tag
	for _, feed := range configFeeds {
		for _, t := range feed["tags"].([]any) {
			if t.(string) == tag {
				feedurls = append(feedurls, feed["rssurl"])
			}
		}
	}
	q := tagFeedsQuery(feedurls)
	feedurls = append(feedurls, feedurls...)
	feedurls = append(feedurls, feedurls...)
	rows, err := db.Query(q, feedurls...)
	if err != nil {
		return tf, err
	}
	defer rows.Close()
	for rows.Next() {
		f := new(Feed)
		rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		tf.Feeds = append(tf.Feeds, f)
	}
	return tf, err
}

func QueryTags(configFeeds []map[string]any) (*TagsPage, error) {
	tp := new(TagsPage)
	tp.TagFeedCount = make(map[string]int)
	tp.Feeds = configFeeds
	for _, feed := range tp.Feeds {
		for _, tag := range feed["tags"].([]any) {
			tp.TagFeedCount[tag.(string)]++
		}
	}
	return tp, nil
}

func anyArticleUnread(db *sql.DB, url ...string) (bool, error) {
	var (
		count   int
		sqlArgs []any
	)
	for _, u := range url {
		sqlArgs = append(sqlArgs, u)
	}
	row := db.QueryRow(articlesUneadQuery(len(url)), sqlArgs...)
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("AnyArticleUnread -> row.Scan: " + fmt.Sprintln(err))
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

type SyncDB struct {
	Unread      int
	FeedUrls    []string
	ArticleUrls []string
}
