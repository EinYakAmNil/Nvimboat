package nvimboat

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func (nb *Nvimboat) QueryArticle(url string) (Article, error) {
	var a = Article{Url: url}
	row := nb.singleRow(articleQuery, a.Url)
	err := row.Scan(&a.Guid, &a.Title, &a.Author, &a.FeedUrl, &a.PubDate, &a.Content, &a.Unread)
	if err != nil {
		return a, err
	}
	return a, nil
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

func (nb *Nvimboat) QueryFeed(feedUrl string) (Feed, error) {
	var (
		f = Feed{RssUrl: feedUrl, UnreadCount: 0}
		a Article
	)
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

func (nb *Nvimboat) QueryFilter(query string, inTags, exTags []string) (Filter, error) {
	var (
		f Filter
		a Article
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
	var (
		tf       TagFeeds
		f        Feed
		feedurls []any
	)
	for _, feed := range nb.ConfigFeeds {
		for _, t := range feed["tags"].([]any) {
			if t.(string) == tag {
				feedurls = append(feedurls, feed["rssurl"])
			}
		}
	}
	q := tagFeedsQuery(feedurls)
	feedurls = append(feedurls, feedurls...)
	feedurls = append(feedurls, feedurls...)
	rows, err := nb.multiRow(q, feedurls...)
	if err != nil {
		return tf, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&f.Title, &f.RssUrl, &f.UnreadCount, &f.ArticleCount)
		tf.Feeds = append(tf.Feeds, f)
	}
	return tf, err
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
