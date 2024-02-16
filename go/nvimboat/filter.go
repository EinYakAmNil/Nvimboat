package nvimboat

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/neovim/go-client/nvim"
)

func (f *Filter) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	cols, err := f.columns(unreadOnly)
	if err != nil {
		return
	}
	for _, col := range cols {
		err = addColumn(nv, buffer, col, separator)
		if err != nil {
			return
		}
	}
	return
}

func (f *Filter) columns(unreadOnly bool) ([][]string, error) {
	dates, err := f.PubDateCol()
	if err != nil {
		return nil, err
	}
	return [][]string{f.PrefixCol(), dates, f.AuthorCol(), f.TitleCol(), f.UrlCol()}, nil
}

func (f *Filter) ChildIdx(article Page) (int, error) {
	for i, a := range f.Articles {
		if a.Url == article.(*Article).Url {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find article in filter.")
}

func (f *Filter) QuerySelf(*sql.DB) (Page, error) {
	return f, nil
}

func (f *Filter) QueryChild(db *sql.DB, articleUrl string) (Page, error) {
	article, err := QueryArticle(db, articleUrl)
	return &article, err
}

func (f *Filter) ToggleUnread(nb Nvimboat, urls ...string) (err error) {
	return nil
}

func (f *Filter) MainPrefix() string {
	ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
	if f.UnreadCount > 0 {
		return "N (" + ratio
	}
	return "  (" + ratio
}

func (f *Filter) PrefixCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Prefix())
	}
	return col
}

func (f *Filter) PubDateCol() ([]string, error) {
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

func (f *Filter) AuthorCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Author)
	}
	return col
}

func (f *Filter) TitleCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Title)
	}
	return col
}

func (f *Filter) UrlCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Url)
	}
	return col
}

func (f *Filter) updateUnreadCount() {
	f.UnreadCount = 0
	for _, a := range f.Articles {
		if a.Unread == 1 {
			f.UnreadCount++
		}
	}
}
