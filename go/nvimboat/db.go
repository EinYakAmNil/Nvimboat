package nvimboat

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(dbpath string) (*sql.DB, error) {
	var err error
	if fileExists(dbpath) {
		d, err := sql.Open("sqlite3", dbpath)
		return d, err
	}
	dbDir := path.Dir(dbpath)
	err = os.MkdirAll(dbDir, os.FileMode(0755))
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	_, err = d.Exec(createDB)
	return d, err
}

func QueryMain(db *sql.DB, configFeeds []map[string]any, configFilters []map[string]any) (*MainMenu, error) {
	var (
		err        error
		tmp_filter Filter
		mainmenu = &MainMenu{ConfigFeeds: configFeeds, ConfigFilters: configFilters}
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
	var f = Feed{RssUrl: feedUrl, UnreadCount: 0}
	row := singleRow(db, feedQuery, f.RssUrl)
	err := row.Scan(&f.Title)
	if err != nil {
		return f, err
	}
	rows, err := multiRow(db, feedArticlesQuery, f.RssUrl)
	if err != nil {
		return f, err
	}
	defer rows.Close()

	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Guid, &a.Title, &a.Author, &a.Url, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
		if err != nil {
			return f, err
		}

		f.Articles = append(f.Articles, a)
	}
	f.ArticleCount = len(f.Articles)
	for _, a := range f.Articles {
		if a.Unread == 1 {
			f.UnreadCount++
		}
	}
	return f, nil
}

func QueryFeeds(db *sql.DB) ([]*Feed, error) {
	var feeds []*Feed
	rows, err := multiRow(db, feedsMinimalQuery)
	if err != nil {
		return feeds, err
	}
	defer rows.Close()

	for rows.Next() {
		f := new(Feed)
		err = rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		if err != nil {
			return feeds, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
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
	rows, err := multiRow(db, q, urls...)
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
	row := singleRow(db, articleQuery, a.Url)
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
	rows, err := multiRow(db, q, feedurls...)
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

func (nb *Nvimboat) anyArticleUnread(url ...string) (bool, error) {
	var (
		count   int
		sqlArgs []any
	)
	for _, u := range url {
		sqlArgs = append(sqlArgs, u)
	}
	row := nb.DB.QueryRow(articlesUneadQuery(len(url)), sqlArgs...)
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("AnyArticleUnread -> row.Scan: " + fmt.Sprintln(err))
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func singleRow(db *sql.DB, query string, args ...any) *sql.Row {
	row := db.QueryRow(query, args...)
	return row
}

func multiRow(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
