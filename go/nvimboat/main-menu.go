package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	MainPageFeed struct {
		rssdb.MainPageFeed
		Tags map[string]bool
	}
	MainMenu struct {
		Feeds   []MainPageFeed
		Filters map[string]*Filter
	}
)

// If the id is a URL then Select() assumes, that a feed is being searched.
// Otherwise the id is matched against the name of a filter.
// Errors if no matching filter is found.
func (mm *MainMenu) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	if len(extracUrls(id)) > 0 {
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
	if filter, ok := mm.Filters[id]; ok {
		return filter, nil
	}
	err = fmt.Errorf(`nvimboat/MainMenu.Select: "%s" is not recognized as an URL or found as a filter name.`, id)
	return
}

func (mm *MainMenu) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	var (
		unreadArticlesRatio []string
		titleCol            []string
		urlCol              []string
	)
	for _, f := range mm.Filters {
		var unreadCount int
		for _, a := range f.Articles {
			if a.Unread == 1 {
				unreadCount++
			}
		}
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(unreadCount, len(f.Articles)))
		titleCol = append(titleCol, f.Name)
		urlCol = append(urlCol, f.ID)
	}
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
	switch f := p.(type) {
	case *Feed:
		childTitle := f.Rssurl
		var (
			section     = len(mm.Feeds)
			searchRange = mm.Feeds
		)
		for range mm.Feeds {
			if childTitle > searchRange[section/2].Feedurl {
				idx += section / 2
				searchRange = searchRange[section/2:]
			} else if childTitle < searchRange[section/2].Feedurl {
				searchRange = searchRange[:section/2]
			} else if childTitle == searchRange[section/2].Feedurl {
				idx += section / 2
				return idx + len(mm.Filters), nil
			}
			section = len(searchRange)
		}
		err = fmt.Errorf(
			`nvimboat/MainMenu.ChildIdx: max iterations. "%s" not found in %v`,
			childTitle,
			prettyStruct(mm.Feeds),
		)
		return
	default:
		pageType := fmt.Sprintf("%T", f)
		err = fmt.Errorf(`nvimboat/MainMenu.ChildIdx: cannot find index of type "%s"`, pageType)
		return -1, err
	}
}

func (mm *MainMenu) Back(*Nvimboat) (int, error) {
	return 0, nil
}
