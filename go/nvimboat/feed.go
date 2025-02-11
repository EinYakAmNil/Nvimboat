package nvimboat

import (
	"encoding/json"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Feed struct {
	rssdb.RssFeed
	Tags     []string
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
		if a.Unread == 0 {
			readStatusCol = append(readStatusCol, " ")
		} else if a.Unread == 1 {
			readStatusCol = append(readStatusCol, "N")

		} else {
			err = fmt.Errorf(`nvimboat/Feed.Render: Bad unread number for "%s" in feed %s: %d\n`,
				a.Url,
				f.Rssurl,
				a.Unread,
			)
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
	feedStruct, _ := json.MarshalIndent(f, "", "	")
	pageStruct, _ := json.MarshalIndent(p, "", "	")
	return -1, fmt.Errorf(`"%v" doesn't contain: "%+v"`, string(feedStruct), string(pageStruct))
}

// TODO: This is very buggy.
// Feeds have to handle going back to either the main menu
// or the tags page where they came from.
// The current implementation is just there to pass the tests
func (f *Feed) Back(nb *Nvimboat) (cursor_x int, err error) {
	if len(nb.Pages.Pages) == 2 {
		err = nb.ShowMain(nb.Nvim, "")
		if err != nil {
			err = fmt.Errorf("nvimboat/Feed.Back: %w\n", err)
			return
		}
	}
	return
}
