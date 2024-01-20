package nvimboat

import (
	"os"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

type (
	Page interface {
		// Render()
		// ElementIdx()
	}
	PageStack struct {
		Pages []Page
		top   Page
	}
	Nvimboat struct {
		Config    map[string]any
		PageStack PageStack
		LogFile   *os.File
		plugin    *nvimPlugin.Plugin
		batch     *nvim.Batch
		buffer    *nvim.Buffer
		window    *nvim.Window
	}
	Main struct {
		Filters []Filter
		Feeds []Feed
	}
	Filter struct {
	}
	Feed struct {
	}
	Article struct {
	}
	TagsPage struct {
	}
	TagFeeds struct {
	}
)

const (
	nvimboatState       = "return package.loaded.nvimboat."
	nvimboatEnable      = nvimboatState + "enable()"
	nvimboatDisable     = nvimboatState + "disable()"
	nvimboatConfig      = nvimboatState + "config"
	nvimboatPage        = nvimboatState + "page"
	nvimboatSetPageType = nvimboatState + "page.set(...)"
)

var Actions = []string{
	"enable",
	"disable",
	"show-main",
	"select",
	"back",
}
