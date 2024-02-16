package nvimboat

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/neovim/go-client/nvim"
)

func (f *Feed) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	dates, err := f.PubDateCol()
	if err != nil {
		return
	}
	cols := [][]string{f.PrefixCol(), dates, f.AuthorCol(), f.TitleCol(), f.UrlCol()}
	for _, col := range cols {
		err = addColumn(nv, buffer, col, separator)
		if err != nil {
			return
		}
	}
	return
}

func (f *Feed) ChildIdx(article Page) (int, error) {
	for i, a := range f.Articles {
		if a.Url == article.(*Article).Url {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find article in feed.")
}

func (f *Feed) QuerySelf(db *sql.DB) (Page, error) {
	newFeed, err := QueryFeed(db, f.RssUrl)
	return &newFeed, err
}

func (f *Feed) QueryChild(db *sql.DB, articleUrl string) (Page, error) {
	article, err := QueryArticle(db, articleUrl)
	return &article, err
}

func (f *Feed) ToggleUnread(nb Nvimboat, urls ...string) (err error) {
	var unreadState int
	hasUnread, err := anyArticleUnread(nb.DBHandler, urls...)
	if hasUnread {
		unreadState = 0
	} else {
		unreadState = 1
	}
	nb.SyncDBchan <- SyncDB{Unread: unreadState, ArticleUrls: urls}
	urlMap := make(map[string]bool)
	for _, url := range urls {
		urlMap[url] = true
	}
	for idx, article := range f.Articles {
		if urlMap[article.Url] {
			f.Articles[idx].Unread = unreadState
		}
	}
	err = f.Render(nb.Nvim, *nb.Buffer, nb.UnreadOnly, nb.Config["separator"].(string))
	return
}

func (f *Feed) MainPrefix() string {
	ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
	if f.UnreadCount > 0 {

		return "N (" + ratio
	}
	return "  (" + ratio
}

func (f *Feed) PrefixCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Prefix())
	}
	return col
}

func (f *Feed) PubDateCol() ([]string, error) {
	var (
		col  []string
		err  error
		date string
	)
	for _, a := range f.Articles {
		date, err = unixToDate(a.PubDate)
		if err != nil {
			return nil, err
		}
		col = append(col, date)
	}
	return col, nil
}

func (f *Feed) AuthorCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Author)
	}
	return col
}

func (f *Feed) TitleCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Title)
	}
	return col
}

func (f *Feed) UrlCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Url)
	}
	return col
}

func (f *Feed) updateUnreadCount() {
	f.UnreadCount = 0
	for _, a := range f.Articles {
		if a.Unread == 1 {
			f.UnreadCount++
		}
	}
}
