package nvimboat

import (
	"fmt"
	"slices"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	MainPageFeed struct {
		rssdb.QueryMainPageRow
		Tags map[string]bool
	}
	MainMenu struct {
		Feeds []MainPageFeed
	}
)

// If the id is a URL then Select() assumes, that a feed is being searched.
// Otherwise the id is matched against the name of a filter.
// Errors if no matching filter is found.
func (mm *MainMenu) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	if len(extracUrls(id)) > 0 {
		p, err = selectFeed(dbh, id)
		if err != nil {
			err = fmt.Errorf("nvimboat/MainMenu.Select: %w\n", err)
			return
		}
		return
	}
	if filter, okFilter := Filters[id]; okFilter {
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
	filterNames := make([]string, 0, len(Filters))
	for name := range Filters {
		filterNames = append(filterNames, name)
	}
	slices.Sort(filterNames)
	var (
		unreadCount int64
		f           *Filter
	)
	for _, filterName := range filterNames {
		f = Filters[filterName]
		for _, a := range f.Articles {
			if a.Unread == 1 {
				unreadCount++
			}
		}
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(unreadCount, int64(len(f.Articles))))
		titleCol = append(titleCol, f.Name)
		urlCol = append(urlCol, f.FilterDescription)
		unreadCount = 0
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
		childTitle := f.Title
		var (
			section     = len(mm.Feeds) / 2
			searchRange = mm.Feeds
		)
		for range len(mm.Feeds) {
			if childTitle > searchRange[section].Title {
				idx += section
				searchRange = searchRange[section:]
			} else if childTitle < searchRange[section].Title {
				searchRange = searchRange[:section]
			} else if childTitle == searchRange[section].Title {
				idx += section
				return idx + len(Filters), nil
			}
			section = len(searchRange) / 2
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

func (mm *MainMenu) Back() (int, error) {
	return 0, nil
}

func (mm *MainMenu) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	return
}
