package nvimboat

import (
	"database/sql"

	"github.com/neovim/go-client/nvim"
)

func (a *Article) Select(nb *Nvimboat, url string) (err error) {
	return
}

func (a *Article) Prefix() string {
	if a.Unread == 1 {
		return "N"
	}
	return " "
}

func (a *Article) Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error) {
	lines, err := a.header()
	if err != nil {
		return
	}
	content, err := renderHTML(a.Content)
	if err != nil {
		return
	}
	lines = append(lines, content...)
	lines = append(lines, "", "# Links")
	lines = append(lines, extracUrls(a.Content)...)

	err = setLines(nv, buffer, lines)
	return
}

func (a *Article) ChildIdx(Page) (int, error) {
	return 0, nil
}

func (a *Article) QuerySelf(db *sql.DB) (Page, error) {
	article, err := QueryArticle(db, a.Url)
	return &article, err
}

func (a *Article) QueryChild(*sql.DB, string) (Page, error) {
	return nil, nil
}

func (a *Article) ToggleUnread(nb *Nvimboat, urls ...string) (err error) {
	nb.Pages.Pop()
	err = nb.Pages.Top().ToggleUnread(nb, urls...)
	if err != nil {
		return
	}
	err = nb.Show(nb.Pages.Top())
	if err != nil {
		return
	}
	pos, err := nb.Pages.Top().ChildIdx(a)
	if err != nil {
		return
	}
	err = nb.Nvim.SetWindowCursor(*nb.Window, [2]int{pos + 1, 1})
	return
}

func (a *Article) Delete(nb *Nvimboat, urls ...string) (err error) {
	nb.Pages.Pop()
	pos, err := nb.Pages.Top().ChildIdx(a)
	if err != nil {
		return
	}
	err = nb.Pages.Top().Delete(nb, urls...)
	if err != nil {
		return
	}
	err = nb.Show(nb.Pages.Top())
	if err != nil {
		return
	}
	err = nb.Nvim.SetWindowCursor(*nb.Window, [2]int{pos + 1, 1})
	if err != nil {
		err = nb.Nvim.SetWindowCursor(*nb.Window, [2]int{pos, 1})
	}
	return
}

func (a *Article) header() (lines []string, err error) {
	date, err := unixToDate(a.PubDate)
	lines = []string{
		"Feed: " + a.FeedUrl,
		"Title: " + a.Title,
		"Author: " + a.Author,
		"Date: " + date,
		"Link: " + a.Url,
		"== Article Begin ==",
	}
	return
}
