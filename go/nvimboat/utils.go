package nvimboat

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

func sortMapKeys(m interface{}) (keyList []string) {
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		keyList = append(keyList, key.Interface().(string))
	}
	sort.Strings(keyList)
	return
}

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
