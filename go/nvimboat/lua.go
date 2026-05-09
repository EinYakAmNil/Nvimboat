package nvimboat

import (
	"errors"
	"fmt"
	"time"
)

const (
	luaPackage    = "return package.loaded.nvimboat"
	luaEnable     = luaPackage + ".enable()"
	luaDisable    = luaPackage + ".disable()"
	luaConfig     = luaPackage + ".config"
	luaFeeds      = luaPackage + ".feeds"
	luaFilters    = luaPackage + ".filters"
	luaPages      = luaPackage + ".pages"
	luaPushPage   = luaPages + ":push(...)"
	luaPopPage    = luaPages + ":pop()"
	luaResetPages = luaPages + ":reset()"
)

func parseConfig(rawConfig map[string]any) (err error) {
	var (
		ok        bool
		cacheTime string
	)
	luaGoMap := map[string]*string{
		"cachePath":   &CachePath,
		"cacheTime":   &cacheTime,
		"dbPath":      &DbPath,
		"linkHandler": &LinkHandler,
		"logPath":     &LogPath,
		"separator":   &ColumnSeparator,
		"userAgent":   &UserAgent,
	}
	for luaName, goVar := range luaGoMap {
		*goVar, ok = rawConfig[luaName].(string)
		if !ok {
			err = errors.Join(err, fmt.Errorf(
				`Lua variable %s must be a string, got: %v -> %T`,
				luaName,
				rawConfig[luaName],
				rawConfig[luaName],
			))
		}
	}
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/parseConfig"))
		return
	}
	CacheTime, err = time.ParseDuration(cacheTime)
	if err != nil {
		err = fmt.Errorf("%w, got: %v", err, cacheTime)
		err = errors.Join(err, errors.New("nvimboat/parseConfig"))
		return
	}
	return
}

func parseFeeds(rawFeeds []map[string]any) (feedConfig map[string]*Feed, tagConfig map[string][]string, err error) {
	var (
		rssurl, t                string
		tagSlice                 []any
		okUrl, okTagSlice, okTag bool
	)
	feedConfig = make(map[string]*Feed)
	tagConfig = make(map[string][]string)
	for _, feed := range rawFeeds {
		if rssurl, okUrl = feed["rssurl"].(string); !okUrl {
			err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type string`, feed["rssurl"])
			return
		}
		if tagSlice, okTagSlice = feed["tags"].([]any); !okTagSlice {
			err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type []any`, feed["tags"])
			return
		}
		feedConfig[rssurl] = new(Feed)
		for _, tag := range tagSlice {
			if t, okTag = tag.(string); !okTag {
				err = fmt.Errorf(`nvimboat/parseFeeds: tag "%v" is not of type string`, tag)
				return
			}
			if feedConfig[rssurl].Tags == nil {
				feedConfig[rssurl].Tags = make(map[string]bool)
			}
			feedConfig[rssurl].Tags[t] = true
			feedConfig[rssurl].Rssurl = rssurl
			tagConfig[t] = append(tagConfig[t], rssurl)
		}
	}
	return
}
