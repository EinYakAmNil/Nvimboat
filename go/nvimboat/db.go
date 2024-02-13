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

func (nb *Nvimboat) QueryMain() (*MainMenu, error) {
	var (
		err        error
		tmp_filter Filter
	)
	mainmenu := new(MainMenu)
	mainmenu.Feeds, err = nb.QueryFeeds()
	if err != nil {
		return mainmenu, err
	}
	mainmenu.Filters, err = nb.parseFilters()
	if err != nil {
		return mainmenu, err
	}
	for i, f := range mainmenu.Filters {
		tmp_filter, err = nb.QueryFilter(f.Query, f.IncludeTags, f.ExcludeTags)
		mainmenu.Filters[i].UnreadCount = tmp_filter.UnreadCount
		mainmenu.Filters[i].ArticleCount = tmp_filter.ArticleCount
		mainmenu.Filters[i].Articles = tmp_filter.Articles
		if err != nil {
			return mainmenu, err
		}
	}
	return mainmenu, err
}

func (nb *Nvimboat) QueryFeed(feedUrl string) (Feed, error) {
	var f = Feed{RssUrl: feedUrl, UnreadCount: 0}
	row := nb.singleRow(feedQuery, f.RssUrl)
	err := row.Scan(&f.Title)
	if err != nil {
		return f, err
	}
	rows, err := nb.multiRow(feedArticlesQuery, f.RssUrl)
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
	for _, a := range f.Articles {
		if a.Unread == 1 {
			f.UnreadCount++
		}
	}
	return f, nil
}

func (nb *Nvimboat) QueryFeeds() ([]*Feed, error) {
	var feeds []*Feed
	rows, err := nb.multiRow(feedsMinimalQuery)
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

func (nb *Nvimboat) QueryFilter(query string, inTags, exTags []string) (Filter, error) {
	var (
		f Filter
	)
	f.Query = query
	f.IncludeTags = inTags
	f.ExcludeTags = exTags
	urls := filterTags(nb.ConfigFeeds, inTags, exTags)
	if len(urls) == 0 {
		return f, nil
	}
	q := articlesFilterQuery(query, len(urls))
	rows, err := nb.multiRow(q, urls...)
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

func (nb *Nvimboat) QueryArticle(url string) (Article, error) {
	var a = Article{Url: url}
	row := nb.singleRow(articleQuery, a.Url)
	err := row.Scan(&a.Guid, &a.Title, &a.Author, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
	if err != nil {
		return a, err
	}
	return a, nil
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
