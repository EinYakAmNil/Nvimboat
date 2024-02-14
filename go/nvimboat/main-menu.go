package nvimboat

import (
	"database/sql"
	"fmt"
	"strconv"
)

func (mm *MainMenu) Render(unreadOnly bool) ([][]string, error) {
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
	return [][]string{prefixCol, titleCol, urlCol}, nil
}

func (mm *MainMenu) SubPageIdx(feed Page) (int, error) {
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

func (mm *MainMenu) QuerySelect(db *sql.DB, id string) (Page, error) {
	switch {
	case id[:4] == "http":
		feed, err := QueryFeed(db, id)
		return &feed, err
	case id[:6] == "query:":
		query, inTags, exTags, err := parseFilterID(id)
		filter, err := QueryFilter(db, mm.ConfigFeeds, query, inTags, exTags)
		return &filter, err
	}
	return nil, fmt.Errorf("Couldn't match ID: %s to anything in the main menu", id)
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
