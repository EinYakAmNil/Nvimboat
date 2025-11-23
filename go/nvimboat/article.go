package nvimboat

import (
	"errors"
	"fmt"
	"slices"

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
		err = errors.Join(err, errors.New("nvimboat/Article.Render"))
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
		err = errors.Join(err, errors.New("nvimboat/Article.Render"))
		return
	}
	lines = append(lines, content...)
	lines = append(lines, "", "# Links")
	lines = append(lines, extracUrls(a.Content)...)

	err = setLines(nv, buf, lines)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Render"))
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
		err = fmt.Errorf("Page stack is less than 2. No parent page possible.")
		err = errors.Join(err, errors.New("nvimboat/Article.Back"))
		return
	}
	cursor_x, err = parentPage.ChildIdx(a)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Back"))
		return
	}
	return cursor_x + 1, nil
}

// Assumption: the article is always in the "read" state.
// It will only ever be made unread by this function.
func (a *Article) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	err = dbh.Queries.SetArticlesUnread(dbh.Ctx, []string{a.Url})
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.ToggleRead"))
		return
	}
	// Update parent page on the state change
	parentPage, err := Pages.Pop()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.ToggleRead"))
		return
	}
	idx, err := parentPage.ChildIdx(a)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.ToggleRead"))
		return
	}
	defer Nvim.SetWindowCursor(*NvWindow, [2]int{idx + 1, 0})
	switch p := parentPage.(type) {
	case *Feed:
		p.Articles[idx].Unread = 1
	case *Filter:
		p.Articles[idx].Unread = 1
	default:
		err = fmt.Errorf(`Unknown parent page type: %T`, p)
		err = errors.Join(err, errors.New("nvimboat/Article.ToggleRead"))
		return
	}
	// Go back
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.ToggleRead"))
		return
	}
	return
}

func (a *Article) NextUnread(dbh rssdb.DbHandle) (err error) {
	parentPage := Pages[len(Pages)-2]
	articleIdx, err := parentPage.ChildIdx(a)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
		return
	}
	switch p := parentPage.(type) {
	case *Feed:
		var newArticle Page
		articleCycle := p.Articles[articleIdx+1:]
		articleCycle = append(articleCycle, p.Articles[:articleIdx]...)
		for _, article := range articleCycle {
			if article.Unread == 1 {
				_, err = Pages.Pop()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				newArticle, err = Pages.Top().Select(dbh, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				err = Pages.Push(newArticle, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				err = Pages.Show()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				return
			}
		}
		Log(`No more unread articles for this feed.`)
	case *Filter:
		var newArticle Page
		articleCycle := p.Articles[articleIdx+1:]
		articleCycle = append(articleCycle, p.Articles[:articleIdx]...)
		for _, article := range articleCycle {
			if article.Unread == 1 {
				_, err = Pages.Pop()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				newArticle, err = Pages.Top().Select(dbh, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				err = Pages.Push(newArticle, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				err = Pages.Show()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
					return
				}
				return
			}
		}
		Log(`No more unread articles for this filter.`)
	default:
		err = fmt.Errorf(
			`Finding next unread article "%s" is not implemented for %T`,
			a.Url,
			p,
		)
		err = errors.Join(err, errors.New("nvimboat/Article.NextUnread"))
		return
	}
	return
}

func (a *Article) PrevUnread(dbh rssdb.DbHandle) (err error) {
	if len(Pages) < 2 {
		err = fmt.Errorf(`Length of page stack is too short: %d`, len(Pages))
		err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
		return
	}
	parentPage := Pages[len(Pages)-2]
	articleIdx, err := parentPage.ChildIdx(a)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
		return
	}
	switch p := parentPage.(type) {
	case *Feed:
		var (
			newArticle Page
		)
		for _, article := range slices.Backward(
			append(p.Articles[articleIdx+1:], p.Articles[:articleIdx]...)) {
			if article.Unread == 1 {
				_, err = Pages.Pop()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				newArticle, err = Pages.Top().Select(dbh, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				err = Pages.Push(newArticle, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				err = Pages.Show()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				return
			}
		}
		Log(`No more unread articles for this feed.`)
	case *Filter:
		var newArticle Page
		articleCycle := slices.Clone(p.Articles[articleIdx+1:])
		articleCycle = append(articleCycle, p.Articles[:articleIdx]...)
		slices.Reverse(articleCycle)
		for _, article := range articleCycle {
			if article.Unread == 1 {
				_, err = Pages.Pop()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				newArticle, err = Pages.Top().Select(dbh, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				err = Pages.Push(newArticle, article.Url)
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				err = Pages.Show()
				if err != nil {
					err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
					return
				}
				return
			}
		}
		Log(`No more unread articles for this filter.`)
	default:
		err = fmt.Errorf(
			`Finding previous unread article "%s" is not implemented for %T`,
			a.Url,
			p,
		)
		err = errors.Join(err, errors.New("nvimboat/Article.PrevUnread"))
		return
	}
	return
}

func (a *Article) Delete(dbh rssdb.DbHandle, ids []string) (err error) {
	err = dbh.Queries.DeleteArticles(dbh.Ctx, ids[:1])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Delete"))
		return
	}
	idx, err := a.Back()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Delete"))
		return
	}
	idx -= 1 // idx is 1-based. Need 0-based
	p, err := Pages.Pop()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Delete"))
		return
	}
	switch f := p.(type) {
	case *Feed:
		f.Articles = append(f.Articles[:idx], f.Articles[idx+1:]...)
	case *Filter:
		f.Articles = append(f.Articles[:idx], f.Articles[idx+1:]...)
	default:
		err = errors.New(`Type of parent page is not Feed/Filter.`)
		err = errors.Join(err, errors.New("nvimboat/Article.Delete"))
		return
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/Back"))
		return
	}
	err = Nvim.SetWindowCursor(*NvWindow, [2]int{max(idx-1, 0), 0})
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Article.Delete"))
		return
	}
	return
}
