package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type TagFeeds struct {
	Name string
	Feeds   []rssdb.QueryTagFeedsRow
}

func (tf *TagFeeds) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	feed, err := dbh.Queries.GetFeed(dbh.Ctx, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/TagFeeds.Select: %w\n", err)
		return
	}
	_ = feed
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
	return
}

func (tf *TagFeeds) Back(nb *Nvimboat) (cursor_x int, err error) {
	return
}
