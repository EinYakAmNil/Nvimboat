package nvimboat

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/neovim/go-client/nvim"
)

func (tp *TagsPage) Select(nb *Nvimboat, id string) (err error) {
	page, err := tp.QueryChild(nb.DBHandler, id)
	if err != nil {
		err = fmt.Errorf("error querying tag '%s' in tags page: %v\n", id, err)
		return
	}
	err = nb.Push(page)
	return
}

func (tp *TagsPage) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	err = setLines(nv, buffer, tp.lines())
	return
}

func (tp *TagsPage) ChildIdx(tagFeeds Page) (int, error) {
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
func (tp *TagsPage) QueryChild(db *sql.DB, tag string) (Page, error) {
	feeds, err := QueryTagFeeds(db, tag, tp.Feeds)
	return &feeds, err
}

func (tp *TagsPage) ToggleUnread(nb *Nvimboat, urls ...string) (err error) {
	return nil
}

func (f *TagsPage) Delete(nb *Nvimboat, urls ...string) (err error) {
	return
}

func (tf *TagFeeds) Select(nb *Nvimboat, url string) (err error) {
	page, err := tf.QueryChild(nb.DBHandler, url)
	if err != nil {
		err = fmt.Errorf("error querying feed '%s' in tag '%s': %v\n", url, tf.Tag, err)
		return
	}
	err = nb.Push(page)
	return
}

func (tf *TagFeeds) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	for _, col := range tf.columns(unreadOnly) {
		err = addColumn(nv, buffer, col, separator)
		if err != nil {
			return
		}
	}
	return
}

func (tf *TagFeeds) ChildIdx(feed Page) (int, error) {
	for i, f := range tf.Feeds {
		if f.RssUrl == feed.(*Feed).RssUrl {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find feed in tag.")
}

func (tf *TagFeeds) QuerySelf(db *sql.DB) (Page, error) {
	var feeds []*Feed
	for _, feed := range tf.Feeds {
		f, err := QueryFeed(db, feed.RssUrl)
		if err != nil {
			return tf, err
		}
		feeds = append(feeds, &f)
	}
	tf.Feeds = feeds
	return tf, nil
}

func (tf *TagFeeds) QueryChild(db *sql.DB, feedUrl string) (Page, error) {
	feed, err := QueryFeed(db, feedUrl)
	return &feed, err
}

func (tf *TagFeeds) ToggleUnread(nb *Nvimboat, urls ...string) (err error) {
	return nil
}

func (f *TagFeeds) Delete(nb *Nvimboat, urls ...string) (err error) {
	return
}

func (tp *TagsPage) lines() (lines []string) {
	for tag, count := range tp.TagFeedCount {
		lines = append(lines, fmt.Sprintf("%s (%d)", tag, count))
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	return
}

func (tf *TagFeeds) columns(unreadOnly bool) (columns [][]string) {
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
	columns = [][]string{prefixCol, titleCol, urlCol}
	return
}
