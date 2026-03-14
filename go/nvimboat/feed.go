package nvimboat

import (
	"errors"
	"fmt"
	"slices"

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
		err = errors.Join(err, errors.New("nvimboat/Feed.Select"))
		return
	}
	err = dbh.Queries.SetArticlesRead(dbh.Ctx, []string{id})
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Select"))
		return
	}
	p = &Article{articleInfo}
	idx, err := f.ChildIdx(p)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Select"))
		return
	}
	f.Articles[idx].Unread = 0
	return
}

func (f *Feed) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(f.Articles) == 0 {
		err = setLines(nv, buf, []string{"No Articles found."})
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Render"))
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
			err = fmt.Errorf(`Bad unread number for "%s" in feed %s: %d`,
				a.Url,
				f.Rssurl,
				a.Unread,
			)
			err = errors.Join(err, errors.New("nvimboat/Feed.Render"))
			return
		}
		parsedTime, err = unixToDate(a.Pubdate)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Render"))
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
			err = errors.Join(err, errors.New("nvimboat/Feed.Render"))
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
	err = fmt.Errorf(
		`"%v" doesn't contain: "%+v"`,
		prettyStruct(f),
		prettyStruct(p),
	)
	err = errors.Join(err, errors.New("nvimboat/Feed.ChildIdx"))
	return -1, err
}

func (f *Feed) Back() (cursor_x int, err error) {
	var parentPage Page
	if len(Pages) >= 2 {
		parentPage = Pages[len(Pages)-2]
	} else {
		err = fmt.Errorf(`Page stack is less than 2. No parent page possible.`)
		err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
		return -1, err
	}
	switch pp := parentPage.(type) {
	case *MainMenu:
		dbh, dbErr := rssdb.ConnectDb(DbPath)
		if dbErr != nil {
			dbErr = errors.Join(dbErr, errors.New("nvimboat/Feed.Back"))
			return -1, dbErr
		}
		pp.Feeds, err = dbh.Queries.QueryMainPage(dbh.Ctx)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
			return -1, err
		}
		err = updateFilters(dbh)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
			return -1, err
		}
		cursor_x, err = pp.ChildIdx(f)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
			return -1, err
		}
		return cursor_x + 1, nil
	case *TagFeeds:
		cursor_x, err = pp.ChildIdx(f)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
			return
		}
		return cursor_x + 1, nil
	default:
		pageType := fmt.Sprintf("%T", parentPage)
		err = fmt.Errorf("parent page type is unaccounted for: %s", pageType)
		return -1, err
	}
}

// If any selected articles are unread, then they will be set to read.
// Set all to unread if all selected articles are read.
func (f *Feed) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
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
			err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
			return
		}

	} else {
		err = dbh.Queries.SetArticlesUnread(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
			return
		}
	}
	f.Articles, err = dbh.Queries.GetFeedPage(dbh.Ctx, f.Rssurl)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
		return
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
		return
	}
	return
}

func (f *Feed) NextUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	cursorRow := cursorPosition[0]
	if len(f.Articles) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this feed's article range: %d.`,
			cursorRow,
			len(f.Articles),
		)
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
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
				err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
				return
			}
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles in this feed.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	return
}

func (f *Feed) PrevUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	cursorRow := cursorPosition[0] - 1
	if len(f.Articles) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this feed's article range: %d.`,
			cursorRow,
			len(f.Articles),
		)
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
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
				err = errors.Join(err, errors.New("nvimboat/Feed.PrevUnread"))
				return
			}
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles in this feed.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	return
}

func (f *Feed) Delete(dbh rssdb.DbHandle, ids []string) (err error) {
	err = dbh.Queries.DeleteArticles(dbh.Ctx, ids)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Delete"))
		return
	}
	f, err = selectFeed(dbh, f.Rssurl)
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
