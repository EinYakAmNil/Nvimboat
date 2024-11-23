package nvimboat

import (
	"fmt"
	"strings"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type MainMenu struct {
	Feeds []rssdb.MainPageFeed
}

func (mm *MainMenu) Select(nb *Nvimboat, id string) (err error) {
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		err = fmt.Errorf("Select: %w", err)
		return
	}
	defer trimTrail(nb.Nvim, *nb.Buffer)
	dbh, err := rssdb.ConnectDb(nb.DbPath)
	if err != nil {
		err = fmt.Errorf("Select: %w", err)
		return
	}
	feedPage := new(Feed)
	feedPage.Articles, err = dbh.Queries.GetFeedPage(dbh.Ctx, id)
	if err != nil {
		err = fmt.Errorf("Select: %w", err)
		return
	}
	err = feedPage.Render(nb.Nvim, *nb.Buffer)
	if err != nil {
		err = fmt.Errorf("Select: %w", err)
		return
	}
	pageType := fmt.Sprintf("%T", feedPage)
	_, pageType, _ = strings.Cut(pageType, "nvimboat.")
	err = nb.Nvim.ExecLua(luaPushPage, new(any), pageType, id)
	if err != nil {
		err = fmt.Errorf("MainMenu.Select: %w", err)
		return
	}
	nb.Pages.Push(feedPage)
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
