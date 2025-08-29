package nvimboat

import (
	"fmt"
	"log"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Feed struct {
	rssdb.RssFeed
	Tags     map[string]bool
	Articles []rssdb.GetFeedPageRow
}

func (f *Feed) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	articleInfo, err := dbh.Queries.GetArticle(dbh.Ctx, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/Feed.Select: %w\n", err)
		return
	}
	err = dbh.Queries.SetArticleRead(dbh.Ctx, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/Feed.Select: %w\n", err)
		return
	}
	p = &Article{articleInfo}
	idx, err := f.ChildIdx(p)
	if err != nil {
		err = fmt.Errorf("nvimboat/Feed.Select: %w\n", err)
		return
	}
	f.Articles[idx].Unread = 0
	return
}

func (f *Feed) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(f.Articles) == 0 {
		err = setLines(nv, buf, []string{"No Articles found."})
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Render: %w\n", err)
			return
		}
		return
	}
	var (
		readStatusCol []string
		parsedTime    string
		pubDateCol    []string
		authorCol     []string
		titleCol      []string
		urlCol        []string
	)
	for _, a := range f.Articles {
		switch a.Unread {
		case 0:
			readStatusCol = append(readStatusCol, " ")
		case 1:
			readStatusCol = append(readStatusCol, "N")
		default:
			err = fmt.Errorf(`nvimboat/Feed.Render: Bad unread number for "%s" in feed %s: %d\n`,
				a.Url,
				f.Rssurl,
				a.Unread,
			)
			log.Println(err)
			return
		}
		parsedTime, err = unixToDate(a.Pubdate)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Render: %w\n", err)
			return
		}
		pubDateCol = append(pubDateCol, parsedTime)
		authorCol = append(authorCol, a.Author)
		titleCol = append(titleCol, a.Title)
		urlCol = append(urlCol, a.Url)
	}
	for _, c := range [][]string{readStatusCol, pubDateCol, authorCol, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Render: %w\n", err)
			return
		}
	}
	return
}

func (f *Feed) ChildIdx(p Page) (idx int, err error) {
	childDate := p.(*Article).Pubdate
	var (
		section     = len(f.Articles)
		searchRange = f.Articles
	)
	for range f.Articles {
		if childDate > searchRange[section/2].Pubdate {
			searchRange = searchRange[:section/2]
		} else if childDate < searchRange[section/2].Pubdate {
			idx += section / 2
			searchRange = searchRange[section/2:]
		} else if childDate == searchRange[section/2].Pubdate {
			idx += section / 2
			return
		}
		section = len(searchRange)
	}
	return -1, fmt.Errorf(
		`"%v" doesn't contain: "%+v"`,
		prettyStruct(f),
		prettyStruct(p),
	)
}

// TODO: This is very buggy.
// Feeds have to handle going back to either the main menu
// or the tags page where they came from.
// The current implementation is just there to pass the tests
func (f *Feed) Back(nb *Nvimboat) (cursor_x int, err error) {
	var parentPage Page
	if len(nb.Pages) >= 2 {
		parentPage = nb.Pages[len(nb.Pages)-2]
	} else {
		err = fmt.Errorf("nvimboat/Feed.Back: page stack is less than 2.\nNo parent page possible.\n")
		return -1, err
	}
	switch pp := parentPage.(type) {
	case *MainMenu:
		dbh, dbErr := rssdb.ConnectDb(nb.DbPath)
		if dbErr != nil {
			dbErr = fmt.Errorf("nvimboat/Feed.Back: %w\n", dbErr)
			return -1, dbErr
		}
		mainPageFeeds, err := dbh.Queries.QueryMainPage(dbh.Ctx)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Back: %w\n", err)
			return -1, err
		}
		for idx, feed := range mainPageFeeds {
			pp.Feeds[idx].QueryMainPageRow = feed
		}
		err = pp.UpdateFilters(dbh)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Back: %w\n", err)
			return -1, err
		}
		cursor_x, err = pp.ChildIdx(f)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Back: %w\n", err)
			return -1, err
		}
		return cursor_x + 1, nil
	case *TagFeeds:
		cursor_x, err = pp.ChildIdx(f)
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Back: %w\n", err)
			return
		}
		return cursor_x + 1, nil
	default:
		pageType := fmt.Sprintf("%T", parentPage)
		err = fmt.Errorf("parent page type is unaccounted for: %s", pageType)
		return -1, err
	}
}
