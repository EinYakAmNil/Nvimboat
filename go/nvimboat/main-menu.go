package nvimboat

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/neovim/go-client/nvim"
)

func (mm *MainMenu) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	for _, col := range mm.columns(unreadOnly) {
		err = addColumn(nv, buffer, col, separator)
		if err != nil {
			return
		}
	}
	return
}

func (mm *MainMenu) columns(unreadOnly bool) [][]string {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	for _, f := range mm.Filters {
		prefixCol = append(prefixCol, mainPrefix(f))
		titleCol = append(titleCol, f.Name)
		urlCol = append(urlCol, f.FilterID)
	}
	for _, f := range mm.Feeds {
		prefixCol = append(prefixCol, mainPrefix(f))
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}
}

func (mm *MainMenu) ChildIdx(feed Page) (int, error) {
	switch feed.(type) {
	case *Filter:
		for i, f := range mm.Filters {
			if feed.(*Filter).FilterID == f.FilterID {
				return i, nil
			}
		}
	case *Feed:
		for i, f := range mm.Feeds {
			if feed.(*Feed).RssUrl == f.RssUrl {
				return i + len(mm.Filters), nil
			}
		}
	default:
		return 0, nil
	}
	return 0, fmt.Errorf("Couldn't find feed/filter.")
}

func (mm *MainMenu) QuerySelf(db *sql.DB) (Page, error) {
	mainmenu, err := QueryMain(db, mm.ConfigFeeds, mm.ConfigFilters)
	return mainmenu, err
}

func (mm *MainMenu) QueryChild(db *sql.DB, id string) (Page, error) {
	switch {
	case id[:4] == "http":
		feed, err := QueryFeed(db, id)
		return &feed, err
	case id[:6] == "query:":
		query, inTags, exTags := parseFilterID(id)
		filter, err := QueryFilter(db, mm.ConfigFeeds, query, inTags, exTags)
		filter.FilterID = id
		return &filter, err
	}
	return nil, fmt.Errorf("Couldn't match ID: %s to anything in the main menu", id)
}

func (mm *MainMenu) ToggleUnread(nb *Nvimboat, urls ...string) (err error) {
	return nil
}

func mainPrefix(p Page) string {
	switch f := p.(type) {
	case *Filter:
		ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
		if f.UnreadCount > 0 {
			return "N (" + ratio
		}
		return "  (" + ratio
	case *Feed:
		ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
		if f.UnreadCount > 0 {
			return "N (" + ratio
		}
		return "  (" + ratio
	}
	return ""
}
