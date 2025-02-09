package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type MainMenu struct {
	Feeds []rssdb.MainPageFeed
}

func (mm *MainMenu) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	feed := new(Feed)
	feed.Articles, err = dbh.Queries.GetFeedPage(dbh.Ctx, id)
	feed.Rssurl = id
	if err != nil {
		err = fmt.Errorf("Select: %w", err)
		return
	}
	p = feed
	return
}

func (mm *MainMenu) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	var (
		unreadArticlesRatio []string
		titleCol            []string
		urlCol              []string
	)
	for _, f := range mm.Feeds {
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(f.UnreadCount, f.ArticleCount))
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.Feedurl)
	}
	for _, c := range [][]string{unreadArticlesRatio, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = fmt.Errorf("MainMenu.Render: %w", err)
			return
		}
	}
	return
}

func (mm *MainMenu) ChildIdx(p Page) (idx int, err error) {
	childTitle := p.(*Feed).Title
	var (
		section     = len(mm.Feeds)
		searchRange = mm.Feeds
	)
	for range mm.Feeds {
		if childTitle > searchRange[section/2].Title {
			idx += section / 2
			searchRange = searchRange[section/2:]
		} else if childTitle < searchRange[section/2].Title {
			searchRange = searchRange[:section/2]
		} else if childTitle == searchRange[section/2].Title {
			idx += section / 2
			return
		}
		section = len(searchRange)
	}
	panic("max iterations!")
}

func (mm *MainMenu) Back(*Nvimboat) (int, error) {
	return 0, nil
}
