package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

type (
	Page interface {
		Render() ([][]string, error)
		ElementIdx(Page) (int, error)
	}
	PageStack struct {
		Pages []Page
		top   Page
	}
	Nvimboat struct {
		Config        map[string]any
		PageStack     PageStack
		ConfigFeeds   []map[string]any
		ConfigFilters []map[string]any
		LogFile       *os.File
		DB            *sql.DB
		plugin        *nvimPlugin.Plugin
		batch         *nvim.Batch
		buffer        *nvim.Buffer
		window        *nvim.Window
	}
	MainMenu struct {
		Filters []Filter
		Feeds   []Feed
	}
	Filter struct {
		Name         string
		FilterID     string
		Query        string
		IncludeTags  []string
		ExcludeTags  []string
		UnreadCount  int
		ArticleCount int
		Articles     []Article
	}
	Feed struct {
		Title        string
		RssUrl       string
		UnreadCount  int
		ArticleCount int
		Articles     []Article
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
		Feeds []Feed
	}
)

const (
	nvimboatState       = "return package.loaded.nvimboat."
	nvimboatEnable      = nvimboatState + "enable()"
	nvimboatDisable     = nvimboatState + "disable()"
	nvimboatConfig      = nvimboatState + "config"
	nvimboatFeeds       = nvimboatState + "feeds"
	nvimboatFilters     = nvimboatState + "filters"
	nvimboatPage        = nvimboatState + "page"
	nvimboatSetPageType = nvimboatState + "page.set(...)"
)

var Actions = []string{
	"enable",
	"disable",
	"show-main",
	"show-tags",
	"select",
	"back",
}
