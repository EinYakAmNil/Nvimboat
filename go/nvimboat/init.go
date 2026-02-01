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
	FeedConfig      map[string][]string
	Feeds           []*Feed
	FilterConfig    []*Filter
	Filters         map[string]*Filter
	LinkHandler     string
	LogPath         string
	ColumnSeparator string
	NvBuffer        *nvim.Buffer
	NvWindow        *nvim.Window
	Nvim            *nvim.Nvim
	Pages           PageStack
)

func initNvimboat(nv *nvim.Nvim) (err error) {
	rawConfig := make(map[string]any)
	rawFeeds := new([]map[string]any)
	rawFilters := new([]map[string]any)
	Filters = make(map[string]*Filter)
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
	FeedConfig, err = parseFeeds(*rawFeeds)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	feedConfig, err := parseFeeds(*rawFeeds)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	for feedurl, tags := range feedConfig {
		f := new(Feed)
		t := make(map[string]bool)
		for _, tag := range tags {
			t[tag] = true
		}
		f.Rssurl = feedurl
		f.Tags = t
		Feeds = append(Feeds, f)
	}
	for _, rawFilter := range *rawFilters {
		filter := new(Filter)
		*filter, err = parseFilter(rawFilter)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
			return
		}
		Filters[filter.Name] = filter
	}
	err = Nvim.SetWindowOption(*NvWindow, "wrap", false)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/initNvimboat"))
		return
	}
	return
}
