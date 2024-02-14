package nvimboat

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
)

func (tp *TagsPage) Render(unreadOnly bool) ([][]string, error) {
	var (
		lines  []string
		prefix string
	)
	for tag, count := range tp.TagFeedCount {
		lines = append(lines, fmt.Sprintf("%s %s (%d)", prefix, tag, count))
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	return [][]string{lines}, nil
}

func (tp *TagsPage) SubPageIdx(tagFeeds Page) (int, error) {
	var tags []string
	for tag := range tp.TagFeedCount {
		tags = append(tags, tag)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	for idx, tag := range tags {
		if tag == tagFeeds.(*TagFeeds).Tag {
			return idx, nil
		}
	}
	return 0, nil
}

func (tp *TagsPage) QuerySelf(*sql.DB) (Page, error) {
	return QueryTags(tp.Feeds)
}
func (tp *TagsPage) QuerySelect(db *sql.DB, tag string) (Page, error) {
	feeds, err := QueryTagFeeds(db, tag, tp.Feeds)
	return &feeds, err
}

func (tf *TagFeeds) Render(unreadOnly bool) ([][]string, error) {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	for _, f := range tf.Feeds {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}, nil
}

func (tf *TagFeeds) SubPageIdx(feed Page) (int, error) {
	for i, f := range tf.Feeds {
		if f.RssUrl == feed.(*Feed).RssUrl {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find feed in tag.")
}

func (tf *TagFeeds) QuerySelf(db *sql.DB) (Page, error) {
	var (
		feeds []*Feed
		f     Feed
		err   error
	)
	for _, feed := range tf.Feeds {
		f, err = QueryFeed(db, feed.RssUrl)
		feeds = append(feeds, &f)
	}
	tf.Feeds = feeds
	return tf, err
}

func (tf *TagFeeds) QuerySelect(*sql.DB, string) (Page, error) {
	return nil, nil
}
