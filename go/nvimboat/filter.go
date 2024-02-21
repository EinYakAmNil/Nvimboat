package nvimboat

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/neovim/go-client/nvim"
)

func (f *Filter) Select(nb *Nvimboat, url string) (err error) {
	article, err := f.QueryChild(nb.DBHandler, url)
	if err != nil {
		err = fmt.Errorf("error querying article '%s' in feed: %v\n", url, err)
		return
	}
	err = nb.Push(article)
	if err != nil {
		err = fmt.Errorf("error pushing '%+v' on page stack in %+v: %v\n", article, f, err)
	}
	if a, ok := article.(*Article); ok {
		err = f.SetArticleRead(*nb, *a)
		return
	}
	err = fmt.Errorf("%+v is not an article\n", article)
	return
}

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
	return 0, fmt.Errorf("Couldn't find article in filter.")
}

func (f *Filter) QuerySelf(*sql.DB) (Page, error) {
	return f, nil
}

func (f *Filter) QueryChild(db *sql.DB, articleUrl string) (Page, error) {
	for _, article := range f.Articles {
		if article.Url == articleUrl {
			return article, nil
		}
	}
	article, err := QueryArticle(db, articleUrl)
	return &article, err
}

func (f *Filter) ToggleUnread(nb *Nvimboat, urls ...string) (err error) {
	var unreadState int
	hasUnread, err := anyArticleUnread(nb.DBHandler, urls...)
	if err != nil {
		err = fmt.Errorf("error querying article unread state for %v: %v\n", urls, err)
		return
	}
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
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		return
	}
	err = f.Render(nb.Nvim, *nb.Buffer, nb.UnreadOnly, nb.Config["separator"].(string))
	return
}

func (f *Filter) FindUnread(direction string, a Article) (article Article, err error) {
	idx, err := f.ChildIdx(&a)
	if err != nil {
		return
	}
	switch direction {
	case "next":
		for i := idx + 1; i < len(f.Articles); i++ {
			if f.Articles[i].Unread == 1 {
				article = *f.Articles[i]
				return
			}
		}
		return a, nil
	case "prev":
		for i := idx - 1; i >= 0; i-- {
			if f.Articles[i].Unread == 1 {
				article = *f.Articles[i]
				return
			}
		}
		return a, nil
	default:
		return a, fmt.Errorf("unknown direction: %s\n", direction)
	}
}

func (f *Filter) SetArticleRead(nb Nvimboat, article Article) (err error) {
	nb.SyncDBchan <- SyncDB{Unread: 0, ArticleUrls: []string{article.Url}}
	idx, err := f.ChildIdx(&article)
	if err != nil {
		return
	}
	f.Articles[idx].Unread = 0
	return
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
