package nvimboat

import (
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	Nvimboat struct {
		CachePath    string
		CacheTime    time.Duration
		FeedConfig   map[string][]string
		FilterConfig []*Filter
		LinkHandler  string
		LogPath      string
		Nvim         *nvim.Nvim
		Window       *nvim.Window
	}
	NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
	Page           interface {
		Select(dbh rssdb.DbHandle, id string) (p Page, err error)
		Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
		ChildIdx(p Page) (idx int, err error)
		Back() (cursor_x int, err error)
		ToggleRead(dbh rssdb.DbHandle, ids []string) (err error)
	}
	PageStack []Page
)
