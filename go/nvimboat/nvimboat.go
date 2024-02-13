package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func (nb *Nvimboat) Prepare(p *nvimPlugin.Plugin) {
	nb.Nvim.Plugin = p
	nb.Config = make(map[string]any)
	nb.Nvim.Batch = p.Nvim.NewBatch()
	nb.Nvim.Buffer = new(nvim.Buffer)
	nb.Nvim.Window = new(nvim.Window)
	nb.Nvim.Batch.CurrentBuffer(nb.Nvim.Buffer)
	nb.Nvim.Batch.CurrentWindow(nb.Nvim.Window)
	nb.Nvim.Batch.ExecLua(nvimboatConfig, &nb.Config)
	nb.Nvim.Batch.ExecLua(nvimboatFeeds, &nb.ConfigFeeds)
	nb.Nvim.Batch.ExecLua(nvimboatFilters, &nb.ConfigFilters)
	nb.Nvim.Batch.SetBufferOption(*nb.Nvim.Buffer, "filetype", "nvimboat")
	nb.Nvim.Batch.SetBufferOption(*nb.Nvim.Buffer, "buftype", "nofile")
	nb.Nvim.Batch.SetWindowOption(*nb.Nvim.Window, "wrap", false)
}

func (nb *Nvimboat) init() error {
	err := nb.Nvim.Batch.Execute()
	if nb.LogFile == nil {
		nb.setupLogging()
	}
	if nb.DB == nil {
		dbpath := nb.Config["dbpath"].(string)
		nb.DB, err = initDB(dbpath)
		if err != nil {
			nb.Log("Error opening the database:")
			nb.Log(err)
		}
	}
	return err
}

type (
	Nvimboat struct {
		Config        map[string]any
		Pages         PageStack
		ConfigFeeds   []map[string]any
		ConfigFilters []map[string]any
		LogFile       *os.File
		DB            *sql.DB
		ChanExecDB    chan DBsync
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
	DBsync struct {
		Unread      int
		FeedUrls    []string
		ArticleUrls []string
	}
	nvimConn struct {
		Plugin *nvimPlugin.Plugin
		Batch  *nvim.Batch
		Buffer *nvim.Buffer
		Window *nvim.Window
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
