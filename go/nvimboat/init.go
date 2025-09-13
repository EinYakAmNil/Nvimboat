package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

var (
	DbPath   string
	Feeds    []*Feed
	Filters  map[string]*Filter
	NvBuffer *nvim.Buffer
	Nvim     *nvim.Nvim
	NvWindow *nvim.Window
	Pages    PageStack
)

func (nb *Nvimboat) init(nv *nvim.Nvim) (err error) {
	rawConfig := make(map[string]any)
	rawFeeds := new([]map[string]any)
	rawFilters := new([]map[string]any)
	Filters = make(map[string]*Filter)
	nb.Nvim = nv
	Nvim = nv
	NvBuffer = new(nvim.Buffer)
	NvWindow = new(nvim.Window)
	execBatch := nv.NewBatch()
	execBatch.CurrentWindow(nb.Window)
	execBatch.CurrentBuffer(NvBuffer)
	execBatch.ExecLua(luaConfig, &rawConfig)
	execBatch.ExecLua(luaFeeds, rawFeeds)
	execBatch.ExecLua(luaFilters, rawFilters)
	err = execBatch.Execute()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.init: %w", err)
		return
	}
	err = parseConfig(nb, rawConfig)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.init parse lua config: %w", err)
		return
	}
	err = SetupLogging(nb.LogPath)
	if err != nil {
		err = fmt.Errorf("Nvimboat init logging: %w", err)
		return
	}
	nb.FeedConfig, err = parseFeeds(*rawFeeds)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.init: %w\n", err)
		return
	}
	feedConfig, err := parseFeeds(*rawFeeds)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.init: %w\n", err)
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
			err = fmt.Errorf("nvimboat/Nvimboat.init: %w\n", err)
			return
		}

		Filters[filter.Name] = filter
	}
	return
}
