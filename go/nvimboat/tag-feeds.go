package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type TagFeeds struct {
	Name  string
	Feeds []rssdb.QueryTagFeedsRow
}

func (tf *TagFeeds) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	p, err = selectFeed(dbh, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/TagFeeds.Select: %w\n", err)
		return
	}
	return
}

func (tf *TagFeeds) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	var (
		unreadArticlesRatio []string
		titleCol            []string
		urlCol              []string
	)
	for _, f := range tf.Feeds {
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(f.UnreadCount, f.ArticleCount))
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.Feedurl)
	}
	for _, c := range [][]string{unreadArticlesRatio, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = fmt.Errorf("nvimboat/TagFeeds.Render: %w\n", err)
			return
		}
	}
	return
}

func (tf *TagFeeds) ChildIdx(p Page) (idx int, err error) {
	switch f := p.(type) {
	case *Feed:
		childTitle := f.Title
		var (
			section     = len(tf.Feeds) / 2
			searchRange = tf.Feeds
		)
		for range len(tf.Feeds) {
			if childTitle > searchRange[section].Title {
				idx += section
				searchRange = searchRange[section:]
			} else if childTitle < searchRange[section].Title {
				searchRange = searchRange[:section]
			} else if childTitle == searchRange[section].Title {
				idx += section
				return idx, nil
			}
			section = len(searchRange) / 2
		}
		err = fmt.Errorf(
			`nvimboat/TagFeeds.ChildIdx: max iterations. "%s" not found in %v`,
			childTitle,
			prettyStruct(tf.Feeds),
		)
	}
	return
}

func (tf *TagFeeds) Back() (cursor_x int, err error) {
	var parentPage Page
	if len(Pages) >= 2 {
		parentPage = Pages[len(Pages)-2]
	} else {
		err = fmt.Errorf("nvimboat/Article.Back: page stack is less than 2.\nNo parent page possible.\n")
		return
	}
	cursor_x, err = parentPage.ChildIdx(tf)
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.Back: %w\n", err)
		return
	}
	return cursor_x + 1, nil
}

func (tf *TagFeeds) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	setFeedsRead := false
	urls := []string{}
checkAnyUnread:
	for _, feed := range tf.Feeds {
		for _, id := range ids {
			if feed.Feedurl == id {
				urls = append(urls, feed.Feedurl)
				if feed.UnreadCount > 0 {
					setFeedsRead = true
				}
				continue checkAnyUnread
			}
		}
	}
	if setFeedsRead {
		err = dbh.Queries.SetFeedsRead(dbh.Ctx, ids)
		if err != nil {
			err = fmt.Errorf("nvimboat/TagFeeds.ToggleRead: %w\n", err)
			return
		}
	} else {
		err = dbh.Queries.SetFeedsUnread(dbh.Ctx, ids)
		if err != nil {
			err = fmt.Errorf("nvimboat/TagFeeds.ToggleRead: %w\n", err)
			return
		}
	}
	tf.Feeds, err = dbh.Queries.QueryTagFeeds(dbh.Ctx, urls)
	if err != nil {
		err = fmt.Errorf("nvimboat/TagFeeds.ToggleRead: %w\n", err)
		return
	}
	err = Pages.Show()
	if err != nil {
		err = fmt.Errorf("nvimboat/TagFeeds.ToggleRead: %w\n", err)
		return
	}
	return
}

func (tf *TagFeeds) NextUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tf *TagFeeds) PrevUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tf *TagFeeds) NextArticle(dbh rssdb.DbHandle) (err error)          { return }
func (tf *TagFeeds) PrevArticle(dbh rssdb.DbHandle) (err error)          { return }
func (tf *TagFeeds) Delete(dbh rssdb.DbHandle, ids []string) (err error) { return }
