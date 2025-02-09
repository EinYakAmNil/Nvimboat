package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Article struct {
	rssdb.GetArticleRow
}

func (a *Article) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	return
}

func (a *Article) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	date, err := unixToDate(a.Pubdate)
	if err != nil {
		err = fmt.Errorf("Article.Render: %w", err)
		return
	}
	lines := []string{
		"Feed: " + a.Feedurl,
		"Title: " + a.Title,
		"Author: " + a.Author,
		"Date: " + date,
		"Link: " + a.Url,
		"== Article Begin ==",
	}
	content, err := renderHTML(a.Content)
	if err != nil {
		err = fmt.Errorf("Article.Render: %w", err)
		return
	}
	lines = append(lines, content...)
	lines = append(lines, "", "# Links")
	lines = append(lines, extracUrls(a.Content)...)

	err = setLines(nv, buf, lines)
	if err != nil {
		err = fmt.Errorf("Article.Render: %w", err)
		return
	}
	return
}

func (a *Article) ChildIdx(Page) (int, error) {
	return -1, fmt.Errorf(`"nvimboat/Article.ChildIdx" should never be called`)
}

func (a *Article) Back(nb *Nvimboat) (cursor_x int, err error) {
	var parentPage Page
	if len(nb.Pages.Pages) >= 2 {
		parentPage = nb.Pages.Pages[len(nb.Pages.Pages)-2]
	} else {
		err = fmt.Errorf("nvimboat/Article.Back: PageStack is less than 2.\nNo parent page possible.\n")
		return
	}
	cursor_x, err = parentPage.ChildIdx(a)
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.Back: %w\n", err)
		return
	}
	return cursor_x + 1, nil
}
