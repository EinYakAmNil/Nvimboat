package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

type (
	Nvimboat struct {
		Config        map[string]any
		Pages         PageStack
		ConfigFeeds   []map[string]any
		ConfigFilters []map[string]any
		LogFile       *os.File
		DB            *sql.DB
		ExecDB        chan DBsync
		Nvim          nvimConn
	}
	Page interface {
		Render(unreadOnly bool) ([][]string, error)
		SubPageIdx(Page) (int, error)
	}
	PageStack struct {
		Pages []*Page
	}
	MainMenu struct {
		Filters []*Filter
		Feeds   []*Feed
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
	nvimConn struct {
		Plugin *nvimPlugin.Plugin
		Batch  *nvim.Batch
		Buffer *nvim.Buffer
		Window *nvim.Window
	}
	DBsync struct {
		Unread      int
		FeedUrls    []string
		ArticleUrls []string
	}
)
