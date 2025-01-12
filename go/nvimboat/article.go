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
		a.Content,
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

func (a *Article) ChildIdx(p Page) (idx int) {
	return
}
