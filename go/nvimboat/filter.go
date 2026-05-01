package nvimboat

import (
	"errors"
	"fmt"
	"slices"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Filter struct {
	rssdb.QueryFilterParams
	Name              string
	FilterDescription string
	IncludeTags       map[string]bool
	ExcludeTags       map[string]bool
	Articles          []rssdb.QueryFilterRow
}

func (f *Filter) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	articleInfo, err := dbh.Queries.GetArticle(dbh.Ctx, id)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Filter.Select"))
		return
	}
	p = &Article{articleInfo}

	Global.ChanAsync <- Async{
		func(...any) (err error) {
			idx, err := f.ChildIdx(p)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/Filter.Select"))
				return
			}
			err = dbh.Queries.SetArticlesRead(dbh.Ctx, []string{id})
			if err != nil {
				return
			}
			f.Articles[idx].Unread = 0

			feedUrl := f.Articles[idx].Feedurl
			feed, err := selectFeed(dbh, feedUrl)
			if err != nil {
				return
			}
			Feeds[feedUrl] = feed
			return
		}, nil,
	}
	return
}

func (f *Filter) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(f.Articles) == 0 {
		err = setLines(nv, buf, []string{"No Articles found."})
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Filter.Render"))
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
			err = fmt.Errorf(`Bad unread number for "%s" in filter %s: %d`,
				a.Url,
				f.Name,
				a.Unread,
			)
			err = errors.Join(err, errors.New("nvimboat/Filter.Render"))
			return
		}
		parsedTime, err = unixToDate(a.Pubdate)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Filter.Render"))
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
			err = errors.Join(err, errors.New("nvimboat/Filter.Render"))
			return
		}
	}
	return
}

func (f *Filter) ChildIdx(p Page) (idx int, err error) {
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
	err = fmt.Errorf(
		`"%v" doesn't contain: "%+v"`,
		prettyStruct(f),
		prettyStruct(p),
	)
	err = errors.Join(err, errors.New("nvimboat/Filter.ChildIdx"))
	return -1, err
}

func (f *Filter) Back() (cursor_x int, err error) {
	for idx, filter := range FilterConfig {
		if f.Name == filter.Name {
			cursor_x = idx + 1
			return
		}
	}
	err = fmt.Errorf(
		"Can't find index for %s",
		prettyStruct(f),
	)
	err = errors.Join(err, errors.New("nvimboat/Filter.Back"))
	return -1, err
}

func (f *Filter) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	setArticlesRead := false
checkAnyUnread:
	for _, a := range f.Articles {
		for _, id := range ids {
			if a.Url == id && a.Unread == 1 {
				setArticlesRead = true
				break checkAnyUnread
			}
		}
	}
	if setArticlesRead {
		err = dbh.Queries.SetArticlesRead(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Filter.ToggleRead"))
			return
		}
	outer1:
		for i, a := range f.Articles {
			for _, id := range ids {
				if a.Url == id && a.Unread == 1 {
					f.Articles[i].Unread = 0
					continue outer1
				}
			}
		}
	} else {
		err = dbh.Queries.SetArticlesUnread(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Filter.ToggleRead"))
			return
		}
	outer2:
		for i, a := range f.Articles {
			for _, id := range ids {
				if a.Url == id && a.Unread == 0 {
					f.Articles[i].Unread = 1
					continue outer2
				}
			}
		}
	}
	return
}

func (f *Filter) NextUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	cursorRow := cursorPosition[0]
	if len(f.Articles) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this filter's article range: %d.`,
			cursorRow,
			len(f.Articles),
		)
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	for i, article := range append(f.Articles[cursorRow:], f.Articles[:cursorRow]...) {
		if article.Unread == 1 {
			err = setCursorUnread(
				(i+cursorRow)%len(f.Articles)+1,
				cursorPosition[1],
				len(f.Articles),
				article,
			)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
				return
			}
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles in this filter.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	return
}

func (f *Filter) PrevUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	cursorRow := cursorPosition[0] - 1
	if len(f.Articles) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this filter's article range: %d.`,
			cursorRow,
			len(f.Articles),
		)
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	for i, article := range slices.Backward(
		append(f.Articles[cursorRow:], f.Articles[:cursorRow]...)) {
		if article.Unread == 1 {
			err = setCursorUnread(
				(i+cursorRow)%len(f.Articles)+1,
				cursorPosition[1],
				len(f.Articles),
				article,
			)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/Filter.PrevUnread"))
				return
			}
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles in this filter.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Filter.NextUnread"))
		return
	}
	return
}

func (f *Filter) Delete(dbh rssdb.DbHandle, ids []string) (err error) {
	err = dbh.Queries.DeleteArticles(dbh.Ctx, ids)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Delete"))
		return
	}
	f.Articles, err = dbh.Queries.QueryFilter(dbh.Ctx, f.QueryFilterParams)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Delete"))
		return
	}
	err = setLines(Nvim, *NvBuffer, []string{""})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	err = f.Render(Nvim, *NvBuffer)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Delete"))
		return
	}
	defer trimTrail(Nvim, *NvBuffer)
	return
}
