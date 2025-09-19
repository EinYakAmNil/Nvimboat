package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Article struct {
	rssdb.GetArticleRow
}

// TODO: Call link handler
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

func (a *Article) Back() (cursor_x int, err error) {
	var parentPage Page
	if len(Pages) >= 2 {
		parentPage = Pages[len(Pages)-2]
	} else {
		err = fmt.Errorf("nvimboat/Article.Back: page stack is less than 2.\nNo parent page possible.\n")
		return
	}
	cursor_x, err = parentPage.ChildIdx(a)
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.Back: %w\n", err)
		return
	}
	return cursor_x + 1, nil
}

// Assumption: the article is always in the "read" state.
// It will only ever be made unread by this function.
func (a *Article) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	err = dbh.Queries.SetArticlesUnread(dbh.Ctx, []string{a.Url})
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.ToggleRead: %w\n", err)
		return
	}
	// Update parent page on the state change
	parentPage, err := Pages.Pop()
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.ToggleRead: %w\n", err)
		return
	}
	idx, err := parentPage.ChildIdx(a)
	if err != nil {
		err = fmt.Errorf("nvimboat/Article.Back: %w\n", err)
		return
	}
	defer Nvim.SetWindowCursor(*NvWindow, [2]int{idx + 1, 0})
	switch p := parentPage.(type) {
	case *Feed:
		p.Articles[idx].Unread = 1
	case *Filter:
		p.Articles[idx].Unread = 1
	default:
		err = fmt.Errorf("nvimboat/Article.ToggleRead: Unknown parent page type: %T\n", p)
		return
	}
	// Go back
	err = Pages.Show()
	if err != nil {
		return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
	}
	return
}

func (a *Article) NextUnread(dbh rssdb.DbHandle) (err error)           { return }
func (a *Article) PrevUnread(dbh rssdb.DbHandle) (err error)           { return }
func (a *Article) NextArticle(dbh rssdb.DbHandle) (err error)          { return }
func (a *Article) PrevArticle(dbh rssdb.DbHandle) (err error)          { return }
func (a *Article) Delete(dbh rssdb.DbHandle, ids []string) (err error) { return }
