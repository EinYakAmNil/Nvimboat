package nvimboat

import (
	"errors"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type TagFeeds struct {
	Name string
	Urls []string
}

func (tf *TagFeeds) Feeds() (feeds []*Feed, err error) {
	feedSubset := make(map[string]*Feed)
	for _, url := range tf.Urls {
		if Feeds[url] == nil {
			err = fmt.Errorf(url, `is not in Feeds`)
			err = errors.Join(err, errors.New("nvimboat/TagFeeds.Feeds"))
			return
		}
		feedSubset[url] = Feeds[url]
	}
	feeds = sortFeeds(feedSubset)
	return
}

func (tf *TagFeeds) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	p, err = selectFeed(dbh, id)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.Select"))
		return
	}
	if f, ok := p.(*Feed); ok {
		Feeds[id] = f
		return
	}
	err = fmt.Errorf(`"%s" not of type %T: %T`, id, new(Feed), p)
	err = errors.Join(err, errors.New("nvimboat/TagFeeds.Select"))
	return
}

func (tf *TagFeeds) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	var (
		unreadArticlesRatio []string
		titleCol            []string
		urlCol              []string
	)
	feeds, err := tf.Feeds()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.Render"))
		return
	}
	for _, f := range feeds {
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(f.UnreadCount, f.ArticleCount))
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.Rssurl)
	}
	for _, c := range [][]string{unreadArticlesRatio, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/TagFeeds.Render"))
			return
		}
	}
	return
}

func (tf *TagFeeds) ChildIdx(p Page) (idx int, err error) {
	feeds, err := tf.Feeds()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.ChildIdx"))
		return
	}
	switch f := p.(type) {
	case *Feed:
		childTitle := f.Title
		var (
			section     = len(feeds) / 2
			searchRange = feeds
		)
		for range len(feeds) {
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
			`Max iterations: "%s" not found in %v`,
			childTitle,
			prettyStruct(feeds),
		)
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.ChildIdx"))
	}
	return
}

func (tf *TagFeeds) Back() (cursor_x int, err error) {
	var parentPage Page
	if len(Pages) >= 2 {
		parentPage = Pages[len(Pages)-2]
	} else {
		err = fmt.Errorf(`Page stack is less than 2. No parent page possible.`)
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.Back"))
		return
	}
	cursor_x, err = parentPage.ChildIdx(tf)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.Back"))
		return
	}
	return cursor_x, nil
}

func (tf *TagFeeds) ToggleRead(dbh rssdb.DbHandle, ids []string) (pos [2][2]int, err error) {
	pos, err = getCursorPositions()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.ToggleRead"))
		return
	}
	feeds, err := tf.Feeds()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.ToggleRead"))
		return
	}
	setFeedsRead := false
	urls := []string{}
checkAnyUnread:
	for _, feed := range feeds {
		for _, id := range ids {
			if feed.Rssurl == id {
				urls = append(urls, feed.Rssurl)
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
			err = errors.Join(err, errors.New("nvimboat/TagFeeds.ToggleRead"))
			return
		}
		for _, id := range ids {
			if Feeds[id] == nil {
				err = fmt.Errorf(id, `is not in Feeds`)
				err = errors.Join(err, errors.New("nvimboat/TagFeeds.Feeds"))
				return
			}
			Feeds[id].UnreadCount = 0
		}
	} else {
		err = dbh.Queries.SetFeedsUnread(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/TagFeeds.ToggleRead"))
			return
		}
		for _, id := range ids {
			if Feeds[id] == nil {
				err = fmt.Errorf(id, `is not in Feeds`)
				err = errors.Join(err, errors.New("nvimboat/TagFeeds.Feeds"))
				return
			}
			Feeds[id].UnreadCount = Feeds[id].ArticleCount
		}
	}
	return
}

func (tf *TagFeeds) NextUnread(dbh rssdb.DbHandle) (err error) { return }
func (tf *TagFeeds) PrevUnread(dbh rssdb.DbHandle) (err error) { return }

func (tf *TagFeeds) Delete(dbh rssdb.DbHandle, ids []string) (err error) {
	for _, id := range ids {
		if len(extracUrls(id)) == 0 {
			Log(fmt.Sprintf(`Can't delete articles for "%s".`, id))
			return
		}
	}
	err = dbh.Queries.DeleteFeedArticles(dbh.Ctx, ids)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagFeeds.Delete"))
		return
	}
	for _, id := range ids {
		if Feeds[id] == nil {
			err = fmt.Errorf(id, `is not in Feeds`)
			err = errors.Join(err, errors.New("nvimboat/TagFeeds.Feeds"))
			return
		}
		Feeds[id].ArticleCount = 0
		Feeds[id].UnreadCount = 0
	}
	return
}
