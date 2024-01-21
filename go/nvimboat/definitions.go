package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

type (
	Page interface {
		Render() ([]string, error)
		// ElementIdx()
	}
	PageStack struct {
		Pages []Page
		top   Page
	}
	Nvimboat struct {
		Config      map[string]any
		PageStack   PageStack
		ConfigFeeds []map[string]any
		LogFile     *os.File
		DB          *sql.DB
		plugin      *nvimPlugin.Plugin
		batch       *nvim.Batch
		buffer      *nvim.Buffer
		window      *nvim.Window
	}
	MainMenu struct {
		Filters []Filter
		Feeds   []Feed
	}
	Filter struct {
		Name         string
		Query        string
		IncludeTags  string
		ExcludeTags  string
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
	}
	TagsPage struct {
		Feeds        []map[string]any
		TagFeedCount map[string]int
	}
	TagFeeds struct {
	}
)

const (
	nvimboatState       = "return package.loaded.nvimboat."
	nvimboatEnable      = nvimboatState + "enable()"
	nvimboatDisable     = nvimboatState + "disable()"
	nvimboatConfig      = nvimboatState + "config"
	nvimboatFeeds       = nvimboatState + "feeds"
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
