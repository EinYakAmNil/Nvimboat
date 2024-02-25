package nvimboat

import (
	"database/sql"

	"github.com/neovim/go-client/nvim"
)

type (
	Page interface {
		Select(nb *Nvimboat, url string) (err error)
		Render(nv *nvim.Nvim, buffer nvim.Buffer, unreadOnly bool, separator string) (err error)
		ChildIdx(Page) (int, error)
		QuerySelf(*sql.DB) (Page, error)
		QueryChild(*sql.DB, string) (Page, error)
		ToggleUnread(nb *Nvimboat, urls ...string) (err error)
		Delete(nb *Nvimboat, urls ...string) (err error)
	}
	ArticlesPage interface {
		FindUnread(direction string, a Article) (Article, error)
		SetArticleRead(nb Nvimboat, article Article) error
	}
	PageStack struct {
		Pages []Page
	}
	MainMenu struct {
		ConfigFeeds   []map[string]any
		ConfigFilters []map[string]any
		Filters       []*Filter
		Feeds         []*Feed
	}
	Filter struct {
		Name         string
		FilterID     string
		Query        string
		IncludeTags  []string
		ExcludeTags  []string
		UnreadCount  int
		ArticleCount int
		Articles     []*Article
	}
	Feed struct {
		Title        string
		RssUrl       string
		UnreadCount  int
		ArticleCount int
		Articles     []*Article
	}
	Article struct {
		Author  string
		Content string
		FeedUrl string
		Guid    string
		PubDate int
		Title   string
		Unread  int
		Url     string
	}
	TagsPage struct {
		Feeds        []map[string]any
		TagFeedCount map[string]int
	}
	TagFeeds struct {
		Tag   string
		Feeds []*Feed
	}
)

func (ps *PageStack) Push(newPage Page) {
	ps.Pages = append(ps.Pages, newPage)
}

func (ps *PageStack) Pop() {
	ps.Pages = ps.Pages[:len(ps.Pages)-1]
}

func (ps *PageStack) Top() Page {
	if pageCount := len(ps.Pages); pageCount > 0 {
		return ps.Pages[pageCount-1]
	}
	return nil
}
