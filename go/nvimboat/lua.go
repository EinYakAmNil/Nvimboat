package nvimboat

import (
	"fmt"
	"time"
)

const (
	luaPackage    = "return package.loaded.nvimboat"
	luaEnable     = luaPackage + ".actions.enable()"
	luaDisable    = luaPackage + ".actions.disable()"
	luaConfig     = luaPackage + ".config"
	luaFeeds      = luaPackage + ".feeds"
	luaFilters    = luaPackage + ".filters"
	luaPages      = luaPackage + ".pages"
	luaPushPage   = luaPackage + ".pages:push(...)"
	luaPopPage    = luaPackage + ".pages:pop()"
	luaResetPages = luaPackage + ".pages:reset()"
)

func parseConfig(nb *Nvimboat, rawConfig map[string]any) (err error) {
	logPath, ok := rawConfig["logPath"].(string)
	if !ok {
		err = fmt.Errorf("parseConfig: log path must be a string, got: %v\n", rawConfig["logPath"])
		return
	}
	nb.LogPath = logPath

	cacheTime, ok := rawConfig["cacheTime"].(string)
	if !ok {
		err = fmt.Errorf("parseConfig: cache time must be a string, got: %v\n", rawConfig["cacheTime"])
		return
	}
	ct, err := time.ParseDuration(cacheTime)
	if err != nil {
		err = fmt.Errorf("parseConfig: %w, got: %v", err, cacheTime)
		return
	}
	nb.CacheTime = ct

	cachePath, ok := rawConfig["cachePath"].(string)
	if !ok {
		err = fmt.Errorf("parseConfig: cache path must be a string, got: %v\n", rawConfig["cachePath"])
		return
	}
	nb.CachePath = cachePath

	dbPath, ok := rawConfig["dbPath"].(string)
	if !ok {
		err = fmt.Errorf("parseConfig: database path must be a string, got: %v\n", rawConfig["dbPath"])
		return
	}
	nb.DbPath = dbPath

	linkHandler, ok := rawConfig["linkHandler"].(string)
	if !ok {
		err = fmt.Errorf("parseConfig: link handler must be a string, got: %v\n", rawConfig["linkHandler"])
		return
	}
	nb.LinkHandler = linkHandler
	return
}

func parseFeeds(rawFeeds []map[string]any) (feedConfig map[string][]string, err error) {
	feedConfig = make(map[string][]string)
	for _, feed := range rawFeeds {
		if rssurl, okUrl := feed["rssurl"].(string); okUrl {
			if tags, okTags := feed["tags"].([]any); okTags {
				for i := range tags {
					tag, okTag := tags[i].(string)
					if okTag {
						feedConfig[rssurl] = append(feedConfig[rssurl], tag)
					} else {
						err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type string`, tags[i])
						return
					}
				}
			} else {
				err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type []any`, feed["tags"])
				return
			}
		} else {
			err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type string`, feed["rssurl"])
			return
		}
	}
	return
}

func parseFilters(rawFilters []map[string]any) (filterConfig []*Filter, err error) {
	f := new(Filter)
	for _, filter := range rawFilters {
		if name, okName := filter["name"].(string); okName {
			f.Name = name
		} else {
			err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
			return
		}
		if query, okQuery := filter["query"].(string); okQuery {
			f.Query = query
		} else {
			err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
			return
		}
		if tags, okTags := filter["tags"].([]any); okTags {
			for _, tag := range tags {
				if t, ok := tag.(string); ok {
					if len(t) == 0 {
						err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
						return
					} else if t[0] == '!' {
						f.ExcludeTags = append(f.ExcludeTags, t)
					} else {
						f.IncludeTags = append(f.IncludeTags, t)
					}
				}
			}
		} else {
			err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
			return
		}
		filterConfig = append(filterConfig, f)
	}
	return
}
