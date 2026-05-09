package nvimboat

import (
	"errors"
	"fmt"
	"time"

	"github.com/neovim/go-client/nvim"
)

var (
	CachePath       string
	CacheTime       time.Duration
	DbPath          string
	TagConfig       map[string][]string
	Feeds           map[string]*Feed
	FilterConfig    []*Filter
	LinkHandler     string
	LogPath         string
	ColumnSeparator string
	NvBuffer        *nvim.Buffer
	NvWindow        *nvim.Window
	Nvim            *nvim.Nvim
	Pages           PageStack
	Global          *Nvimboat
	UserAgent       string
)

func initNvimboat(nb *Nvimboat, nv *nvim.Nvim) (err error) {
	Global = nb
	// Global.ChanAsync = make(chan Async)
	rawConfig := make(map[string]any)
	rawFeeds := new([]map[string]any)
	rawFilters := new([]map[string]any)
	Nvim = nv
	NvBuffer = new(nvim.Buffer)
	NvWindow = new(nvim.Window)
	execBatch := nv.NewBatch()
	execBatch.CurrentWindow(NvWindow)
	execBatch.CurrentBuffer(NvBuffer)
	execBatch.ExecLua(luaConfig, &rawConfig)
	execBatch.ExecLua(luaFeeds, rawFeeds)
	execBatch.ExecLua(luaFilters, rawFilters)
	err = execBatch.Execute()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	err = parseConfig(rawConfig)
	if err != nil {
		err = fmt.Errorf(`Parse lua config: %w`, err)
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	err = SetupLogging(LogPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	Feeds, TagConfig, err = parseFeeds(*rawFeeds)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	for _, rawFilter := range *rawFilters {
		filter := new(Filter)
		*filter, err = parseFilter(rawFilter)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
			return
		}
		FilterConfig = append(FilterConfig, filter)
	}
	return
}
