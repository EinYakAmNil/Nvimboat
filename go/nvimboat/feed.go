package nvimboat

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Feed struct {
	rssdb.GetFeedRow
	Tags     map[string]bool
	Articles []rssdb.GetFeedPageRow
}

func (f *Feed) ID() string {
	return f.Rssurl
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
	f.UnreadCount = 0
	for _, a := range f.Articles {
		if a.Unread == 1 {
			f.UnreadCount++
		} else if a.Unread > 1 || a.Unread < 0 {
			err = fmt.Errorf(`Unexpected value for unread: %d`, a.Unread)
			err = errors.Join(err, errors.New("nvimboat/Feed.Select"))
			return
		}
	}
	Global.ChanAsync <- Async{func(...any) (err error) {
		err = updateFilters(dbh)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
			return
		}
		return
	}, nil}
	return
}

func (f *Feed) Open(urls ...string) (err error) {
	cmd := exec.Command(LinkHandler, urls...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	err = cmd.Start()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Open"))
		return
	}
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
	cursor_x, err = parentPage.ChildIdx(f)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.Back"))
		return -1, err
	}
	return cursor_x + 1, nil
}

// If any selected articles are unread, then they will be set to read.
// Set all to unread if all selected articles are read.
func (f *Feed) ToggleRead(dbh rssdb.DbHandle, ids []string) (pos [2][2]int, err error) {
	pos, err = getCursorPositions()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
		return
	}
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
	feed, err := selectFeed(dbh, f.Rssurl)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
		return
	}
	*f = *feed
	Global.ChanAsync <- Async{func(...any) (err error) {
		err = updateFilters(dbh)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Feed.ToggleRead"))
			return
		}
		return
	}, nil}
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
