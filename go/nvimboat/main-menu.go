package nvimboat

import (
	"errors"
	"fmt"
	"slices"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type MainMenu struct {
}

// If the id is a URL then Select() assumes, that a feed is being searched.
// Otherwise the id is matched against the name of a filter.
// Errors if no matching filter is found.
func (mm *MainMenu) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	if len(extracUrls(id)) > 0 {
		p, err = selectFeed(dbh, id)
		if err != nil {
			err = errors.Join(err, errors.New("in nvimboat/MainMenu.Select"))
			return
		}
		if f, ok := p.(*Feed); ok {
			Feeds[id] = f
			return
		}
		err = fmt.Errorf(`"%s" not of type %T: %T`, id, new(Feed), p)
		err = errors.Join(err, errors.New("nvimboat/MainMenu.Select"))
		return
	}
	for _, filter := range FilterConfig {
		if id == filter.Name {
			return filter, nil
		}
	}
	err = fmt.Errorf(
		`"%s" is not recognized as an URL or found as a filter name.`,
		id,
	)
	err = errors.Join(err, errors.New("nvimboat/MainMenu.Select"))
	return
}

func (mm *MainMenu) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	var (
		unreadArticlesRatio []string
		titleCol            []string
		urlCol              []string
		unreadCount         int64
	)
	for _, f := range FilterConfig {
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
	for _, f := range sortFeeds(Feeds) {
		unreadArticlesRatio = append(unreadArticlesRatio, makeUnreadRatio(f.UnreadCount, f.ArticleCount))
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.Rssurl)
	}
	for _, c := range [][]string{unreadArticlesRatio, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/MainMenu.Render"))
			return
		}
	}
	return
}

func (mm *MainMenu) ChildIdx(p Page) (idx int, err error) {
	feedCount := len(Feeds)
	switch f := p.(type) {
	case *Feed:
		childTitle := f.Title
		var (
			section     = feedCount / 2
			searchRange = sortFeeds(Feeds)
		)
		for range feedCount {
			if childTitle > searchRange[section].Title {
				idx += section
				searchRange = searchRange[section:]
			} else if childTitle < searchRange[section].Title {
				searchRange = searchRange[:section]
			} else if childTitle == searchRange[section].Title {
				idx += section
				return idx + len(FilterConfig), nil
			}
			section = len(searchRange) / 2
		}
		err = fmt.Errorf(
			`Max iterations. "%s" not found in %v`,
			childTitle,
			prettyStruct(sortFeeds(Feeds)),
		)
		err = errors.Join(err, errors.New("nvimboat/MainMenu.ChildIdx"))
		return
	default:
		err = fmt.Errorf(`Cannot find index of type "%s"`, fmt.Sprintf("%T", f))
		err = errors.Join(err, errors.New("nvimboat/MainMenu.ChildIdx"))
		return -1, err
	}
}

func (mm *MainMenu) Back() (int, error) {
	return 0, nil
}

func (mm *MainMenu) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	for _, id := range ids {
		if len(extracUrls(id)) == 0 {
			Log(fmt.Sprintf(`Can't toggle read for "%s".`, id))
			return
		}
	}
	setFeedsRead := false
checkAnyUnread:
	for _, f := range sortFeeds(Feeds) {
		for _, id := range ids {
			if f.Rssurl == id && f.UnreadCount > 0 {
				setFeedsRead = true
				break checkAnyUnread
			}
		}
	}
	if setFeedsRead {
		err = dbh.Queries.SetFeedsRead(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/MainMenu.ToggleRead"))
			return
		}
	} else {
		err = dbh.Queries.SetFeedsUnread(dbh.Ctx, ids)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/MainMenu.ToggleRead"))
			return
		}
	}
	mainPageFeeds, err := dbh.Queries.QueryMainPage(dbh.Ctx)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.ToggleRead"))
		return
	}
	for _, mpf := range mainPageFeeds {
		Feeds[mpf.Rssurl].Title = mpf.Title
		Feeds[mpf.Rssurl].Rssurl = mpf.Rssurl
		Feeds[mpf.Rssurl].ArticleCount = mpf.ArticleCount
		Feeds[mpf.Rssurl].UnreadCount = mpf.UnreadCount
	}
	return
}

func (mm *MainMenu) NextUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
		return
	}
	cursorRow := cursorPosition[0]
	var rows []any
	for _, f := range FilterConfig {
		rows = append(rows, f)
	}
	for _, f := range sortFeeds(Feeds) {
		rows = append(rows, f)
	}
	if len(rows) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this filter's article range: %d.`,
			cursorRow,
			len(rows),
		)
		err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
		return
	}
	for i, r := range append(rows[cursorRow:], rows[:cursorRow]...) {
		switch f := r.(type) {
		case *Filter:
			for _, a := range f.Articles {
				if a.Unread == 1 {
					err = setCursorUnread(
						(i+cursorRow)%len(rows)+1,
						cursorPosition[1],
						len(rows),
						f,
					)
					if err != nil {
						err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
						return
					}
					return
				}
			}
		case *Feed:
			if f.UnreadCount > 0 {
				err = setCursorUnread(
					(i+cursorRow)%len(rows)+1,
					cursorPosition[1],
					len(rows),
					f,
				)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
					return
				}
				return
			}
		default:
			err = fmt.Errorf(`Unknown row type "%T" during search for %+v`, f, f)
			err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.NextUnread"))
		return
	}
	return
}

func (mm *MainMenu) PrevUnread(dbh rssdb.DbHandle) (err error) {
	cursorPosition, err := Nvim.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
		return
	}
	cursorRow := cursorPosition[0] - 1
	var rows []any
	for _, f := range FilterConfig {
		rows = append(rows, f)
	}
	for _, f := range sortFeeds(Feeds) {
		rows = append(rows, f)
	}
	if len(rows) < cursorRow {
		err = fmt.Errorf(
			`Cursor row (%d) is outside of this filter's article range: %d.`,
			cursorRow,
			len(rows),
		)
		err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
		return
	}
	for i, r := range slices.Backward(
		append(rows[cursorRow:], rows[:cursorRow]...)) {
		switch f := r.(type) {
		case *Filter:
			for _, a := range f.Articles {
				if a.Unread == 1 {
					err = setCursorUnread(
						(i+cursorRow)%len(rows)+1,
						cursorPosition[1],
						len(rows),
						f,
					)
					if err != nil {
						err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
						return
					}
					return
				}
			}
		case *Feed:
			if f.UnreadCount > 0 {
				err = setCursorUnread(
					(i+cursorRow)%len(rows)+1,
					cursorPosition[1],
					len(rows),
					f,
				)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
					return
				}
				return
			}
		default:
			err = fmt.Errorf(`Unknown row type "%T" during search for %+v`, f, f)
			err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
			return
		}
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "No more unread articles.",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.PrevUnread"))
		return
	}
	return
}

func (mm *MainMenu) Delete(dbh rssdb.DbHandle, ids []string) (err error) {
	for _, id := range ids {
		if len(extracUrls(id)) == 0 {
			Log(fmt.Sprintf(`Can't delete articles for "%s".`, id))
			return
		}
	}
	err = dbh.Queries.DeleteFeedArticles(dbh.Ctx, ids)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.Delete"))
		return
	}
	err = ShowMain(Nvim)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/MainMenu.Delete"))
		return
	}
	return
}
