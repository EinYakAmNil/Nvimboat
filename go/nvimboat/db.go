package nvimboat

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbpath string) (db *sql.DB, err error) {
	if fileExists(dbpath) {
		db, err = sql.Open("sqlite3", dbpath)
		if err != nil {
			err = fmt.Errorf("error opening database '%s': %v\n", dbpath, err)
		}
		return
	}
	dbDir := path.Dir(dbpath)
	err = os.MkdirAll(dbDir, os.FileMode(0755))
	if err != nil {
		err = fmt.Errorf("error creating directory '%s': %v\n", dbDir, err)
		return
	}
	d, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		err = fmt.Errorf("error opening database '%s': %v\n", dbpath, err)
		return
	}
	_, err = d.Exec(createDB)
	if err != nil {
		err = fmt.Errorf("error creating tables for '%s': %v\n", dbpath, err)
	}
	return
}

func QueryMain(db *sql.DB, configFeeds []map[string]any, configFilters []map[string]any) (mainmenu *MainMenu, err error) {
	mainmenu = &MainMenu{ConfigFeeds: configFeeds, ConfigFilters: configFilters}
	mainmenu.Feeds, err = QueryFeeds(db)
	if err != nil {
		err = fmt.Errorf("error querying all feeds: %v\n", err)
		return
	}
	mainmenu.Filters, err = parseFilters(configFilters)
	if err != nil {
		err = fmt.Errorf("error parsing filters '%+v': %v\n", configFilters, err)
		return
	}
	var tmp_filter Filter
	for i, f := range mainmenu.Filters {
		tmp_filter, err = QueryFilter(db, configFeeds, f.Query, f.IncludeTags, f.ExcludeTags)
		if err != nil {
			err = fmt.Errorf("error querying filter '%s': %v\n", mainmenu.Filters[i].FilterID, err)
			return
		}
		mainmenu.Filters[i].UnreadCount = tmp_filter.UnreadCount
		mainmenu.Filters[i].ArticleCount = tmp_filter.ArticleCount
		mainmenu.Filters[i].Articles = tmp_filter.Articles
	}
	return
}

func QueryFeed(db *sql.DB, feedUrl string) (feed Feed, err error) {
	feed = Feed{RssUrl: feedUrl, UnreadCount: 0}
	row := db.QueryRow(feedQuery, feed.RssUrl)
	err = row.Scan(&feed.Title)
	if err != nil {
		err = fmt.Errorf("error querying feed '%s': %v\n", feedUrl, err)
		return
	}
	rows, err := db.Query(feedArticlesQuery, feed.RssUrl)
	if err != nil {
		err = fmt.Errorf("error querying for articles of feed '%s' with query '%s': %v\n", feedUrl, feedArticlesQuery, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Guid, &a.Title, &a.Author, &a.Url, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
		if err != nil {
			err = fmt.Errorf("error assigning values for article '%+v' in feed '%s': %v\n", a, feedUrl, err)
			return
		}
		feed.Articles = append(feed.Articles, a)
	}
	feed.ArticleCount = len(feed.Articles)
	for _, a := range feed.Articles {
		if a.Unread == 1 {
			feed.UnreadCount++
		}
	}
	return
}

func QueryFeeds(db *sql.DB) (feeds []*Feed, err error) {
	rows, err := db.Query(feedsMinimalQuery)
	if err != nil {
		err = fmt.Errorf("error querying all feeds with '%s': %v\n", feedsMinimalQuery, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		f := new(Feed)
		err = rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		if err != nil {
			err = fmt.Errorf("error assigning values for feed '%+v' in QueryFeeds: %v\n", f, err)
			return
		}
		feeds = append(feeds, f)
	}
	return
}

func QueryFilter(db *sql.DB, configFeeds []map[string]any, query string, inTags, exTags []string) (filter Filter, err error) {
	filter.Query = query
	filter.IncludeTags = inTags
	filter.ExcludeTags = exTags
	urls := filterTags(configFeeds, inTags, exTags)
	if len(urls) == 0 {
		return
	}
	q := articlesFilterQuery(query, len(urls))
	rows, err := db.Query(q, urls...)
	if err != nil {
		err = fmt.Errorf("error querying articles for ")
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Guid, &a.Title, &a.Author, &a.Url, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
		if err != nil {
			err = fmt.Errorf("error assigning value to article '%+v' in QueryFilter '%s': %v\n", a, q, err)
			return
		}
		filter.ArticleCount++
		if a.Unread == 1 {
			filter.UnreadCount++
		}
		filter.Articles = append(filter.Articles, a)
	}
	return
}

func QueryArticle(db *sql.DB, url string) (article Article, err error) {
	article = Article{Url: url}
	row := db.QueryRow(articleQuery, article.Url)
	err = row.Scan(&article.Guid, &article.Title, &article.Author, &article.FeedUrl, &article.PubDate, &article.Content, &article.Unread)
	if err != nil {
		err = fmt.Errorf("error querying for article '%+v' with query '%s': %v\n", article, articleQuery, err)
	}
	return
}

func QueryTagFeeds(db *sql.DB, tag string, configFeeds []map[string]any) (tf TagFeeds, err error) {
	tf.Tag = tag

	var feedurls []any
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
		err = fmt.Errorf("error querying feed of tag '%s' with '%s': %v\n", tag, q, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := new(Feed)
		err = rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		if err != nil {
			err = fmt.Errorf("error assigning value to feed '%+v': %v\n", f, err)
			return
		}
		tf.Feeds = append(tf.Feeds, f)
	}
	return
}

func QueryTags(configFeeds []map[string]any) (tp *TagsPage, err error) {
	tp = new(TagsPage)
	tp.TagFeedCount = make(map[string]int)
	tp.Feeds = configFeeds
	for _, feed := range tp.Feeds {
		for _, tag := range feed["tags"].([]any) {
			tp.TagFeedCount[tag.(string)]++
		}
	}
	return
}

func anyArticleUnread(db *sql.DB, url ...string) (hasUnread bool, err error) {
	var (
		count   int
		sqlArgs []any
	)
	for _, u := range url {
		sqlArgs = append(sqlArgs, u)
	}
	row := db.QueryRow(articlesUneadQuery(len(url)), sqlArgs...)
	err = row.Scan(&count)
	if err != nil {
		err = fmt.Errorf("AnyArticleUnread -> row.Scan: " + fmt.Sprintln(err))
		hasUnread = false
		return
	}
	if count > 0 {
		hasUnread = true
		return
	}
	hasUnread = false
	return
}

type SyncDB struct {
	Unread      int
	Delete      bool
	FeedUrls    []string
	ArticleUrls []string
}
