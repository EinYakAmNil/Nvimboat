package nvimboat

import (
	"errors"
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
	}
	for luaName, GoVar := range luaGoMap {
		*GoVar, ok = rawConfig[luaName].(string)
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
