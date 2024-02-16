package nvimboat

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) init(nv *nvim.Nvim) (err error) {
	nb.Nvim = nv
	nb.SyncDBchan = make(chan SyncDB)
	nb.Config = make(map[string]any)
	nb.Window = new(nvim.Window)
	nb.Buffer = new(nvim.Buffer)
	nb.UnreadOnly = false
	execBatch := nv.NewBatch()
	execBatch.CurrentWindow(nb.Window)
	execBatch.CurrentBuffer(nb.Buffer)
	execBatch.ExecLua(nvimboatConfig, &nb.Config)
	execBatch.ExecLua(nvimboatFeeds, &nb.Feeds)
	execBatch.ExecLua(nvimboatFilters, &nb.Filters)
	execBatch.SetBufferOption(*nb.Buffer, "filetype", "nvimboat")
	execBatch.SetBufferOption(*nb.Buffer, "buftype", "nofile")
	execBatch.SetWindowOption(*nb.Window, "wrap", false)
	err = execBatch.Execute()
	if err != nil {
		return
	}
	nb.DBHandler, err = initDB(nb.Config["dbpath"].(string))
	if err != nil {
		return
	}
	err = SetupLogging(nb.Config["log"].(string))
	return
}

func articlesUneadQuery(n int) string {
	if n == 0 {
		return ""
	}
	const (
		prefix = `SELECT COUNT(unread) FROM rss_item WHERE unread = 1 AND url IN (?`
		suffix = `)`
	)
	if n < 2 {
		return prefix + suffix
	}
	articleCount := strings.Repeat(", ?", n-1)

	return prefix + articleCount + suffix
}

func (nb *Nvimboat) setPageType(p Page) (err error) {
	t := pageTypeString(p)
	err = nb.Nvim.ExecLua(nvimboatSetPageType, new(any), t)
	return
}

func addColumn(nv *nvim.Nvim, buffer nvim.Buffer, col []string, separator string) (err error) {
	currentLines, err := nv.BufferLines(buffer, 0, -1, false)
	if err != nil {
		return
	}
	var (
		diff  int
		lines = []string{}
	)
	diff = (len(col) - len(currentLines))
	for i := 0; i < diff; i++ {
		currentLines = append(currentLines, []byte{})
	}
	for i, c := range col {
		lines = append(lines, string(currentLines[i])+separator+c)
	}
	err = setLines(nv, buffer, lines)
	if err != nil {
		return
	}
	vcl, err := virtColLens(nv)
	if err != nil {
		return
	}
	maxLineLen := slices.Max(vcl)

	for i, l := range lines {
		diff = maxLineLen - vcl[i]
		lines[i] = l + strings.Repeat(" ", diff)
	}
	err = setLines(nv, buffer, lines)
	return err
}

func virtColLens(nv *nvim.Nvim) (evalResult []int, err error) {
	virtCols := "map(range(1, line('$')), \"virtcol([v:val, '$'])\")"
	err = nv.Eval(virtCols, &evalResult)
	return
}

func trimTrail(nv *nvim.Nvim, buffer nvim.Buffer) (err error) {
	currentLines, err := nv.BufferLines(buffer, 0, -1, false)
	if err != nil {
		return
	}
	var lines []string
	for _, l := range currentLines {
		lines = append(lines, strings.TrimRight(string(l), " "))
	}
	err = setLines(nv, buffer, lines)
	return
}

func parseFilters(configFilters []map[string]any) (filters []*Filter, err error) {
	for _, filter := range configFilters {
		f := new(Filter)
		if name, ok := filter["name"]; ok {
			f.Name = name.(string)
		} else {
			return filters, fmt.Errorf("Failed to parse: %v", filter)
		}
		if query, ok := filter["query"]; ok {
			f.Query = query.(string)
			f.FilterID = "query: " + query.(string) + ", tags: "
		} else {
			return filters, fmt.Errorf("Failed to parse: %v", filter)
		}
		if tags, ok := filter["tags"]; ok {
			for _, tag := range tags.([]any) {
				if len(tag.(string)) == 0 {
					continue
				}
				if tag.(string)[0] != '!' {
					f.IncludeTags = append(f.IncludeTags, tag.(string))
				} else {
					f.ExcludeTags = append(f.ExcludeTags, tag.(string)[1:])
				}
				f.FilterID += tag.(string) + ", "
			}
		}
		if f.FilterID[len(f.FilterID)-2:] == ", " {
			f.FilterID = f.FilterID[:len(f.FilterID)-2]
		}
		filters = append(filters, f)
	}
	return
}

func parseFilterID(id string) (query string, includeTags []string, excludeTags []string) {
	query, rawTags, _ := strings.Cut(id, ", ")
	_, query, _ = strings.Cut(query, "query: ")
	_, rawTags, _ = strings.Cut(rawTags, "tags: ")
	tags := strings.Split(rawTags, ", ")
	for _, t := range tags {
		if string(t[0]) == "!" {
			excludeTags = append(excludeTags, t[1:])
		} else {
			includeTags = append(includeTags, t)
		}
	}
	return
}

func strings2bytes(stringSlice []string) (byteSlices [][]byte) {
	for _, s := range stringSlice {
		byteSlices = append(byteSlices, []byte(s))
	}
	return
}

func unixToDate(unixTime int) (string, error) {
	tz, err := time.LoadLocation("Local")
	if err != nil {
		return "", err
	}
	ut := time.Unix(int64(unixTime), 0)
	dateString := ut.In(tz).Format("02 Jan 06")

	return dateString, nil
}

func extracUrls(content string) (links []string) {
	re := regexp.MustCompile(`\b((?:https?|ftp|file):\/\/[-a-zA-Z0-9+&@#\/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#\/%=~_|])`)
	matches := re.FindAll([]byte(content), -1)
	for _, l := range matches {
		links = append(links, string(l))
	}
	return links
}

func renderHTML(content string) ([]string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(content)
	if err != nil {
		return nil, err
	}
	return strings.Split(markdown, "\n"), nil
}

func pageTypeString(p Page) (name string) {
	fullName := fmt.Sprintf("%T", p)
	_, name, _ = strings.Cut(fullName, "nvimboat.")
	return
}

func articlesFilterQuery(query string, n int) string {
	const (
		prefix = `
		SELECT guid, title, author, url, feedurl, pubDate, content, unread
		FROM rss_item WHERE deleted = 0 AND feedurl in (?
		`
		suffix = ` ORDER BY pubDate DESC`
	)
	var glue string
	if query != "" {
		glue = `) AND `
	} else {
		glue = `) `
	}
	if n < 2 {
		return prefix + glue + query + suffix
	}
	articleCount := strings.Repeat(", ?", n-1)

	return prefix + articleCount + glue + query + suffix
}

func filterTags(configFeeds []map[string]any, inTags, exTags []string) (urls []any) {
	feedurls := make(map[string]bool)
	for _, feed := range configFeeds {
		if subset(inTags, anyToString(feed["tags"].([]any))) {
			feedurls[feed["rssurl"].(string)] = true
		}
	}
	for _, feed := range configFeeds {
		for _, tag := range anyToString(feed["tags"].([]any)) {
			if slices.Contains(exTags, tag) {
				delete(feedurls, feed["rssurl"].(string))
				continue
			}
		}
	}
	for url := range feedurls {
		urls = append(urls, url)
	}
	return
}

func setLines(nv *nvim.Nvim, buffer nvim.Buffer, lines []string) (err error) {
	err = nv.SetBufferLines(buffer, 0, -1, false, strings2bytes(lines))
	return
}

func tagFeedsQuery(feedurls []any) string {
	n := len(feedurls)
	if n == 0 {
		return ""
	}
	p1 := `
	SELECT rss_feed.title, c.* FROM rss_feed
	LEFT JOIN (
	SELECT a.feedurl, b.unreadCount, a.articleCount
	FROM (
	SELECT feedurl, COUNT(*) AS articleCount
	FROM rss_item WHERE feedurl in (?`
	p2 := `)
	GROUP BY feedurl
	) a
	LEFT JOIN (
	SELECT feedurl, sum(unread) AS unreadCount
	FROM rss_item WHERE feedurl in (?`
	p3 := `)
	GROUP BY feedurl
	) b
	ON a.feedurl = b.feedurl
	) c
	ON rss_feed.rssurl = c.feedurl
	WHERE rssurl in (?`
	p4 := `)
	ORDER BY rss_feed.title`
	if n < 2 {
		return p1 + p2 + p3 + p4
	}
	reps := strings.Repeat(", ?", n-1)
	return p1 + reps + p2 + reps + p3 + reps + p4
}

func subset(first, second []string) bool {
	set := make(map[string]bool)
	for _, value := range second {
		set[value] = true
	}
	for _, value := range first {
		if item, found := set[value]; !found {
			return false
		} else if !item {
			return false
		}
	}
	return true
}

func anyToString(base []any) (conv []string) {
	for _, a := range base {
		conv = append(conv, a.(string))
	}
	return
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}

func sortedMapKeys(m interface{}) (keyList []string) {
	keys := reflect.ValueOf(m).MapKeys()

	for _, key := range keys {
		keyList = append(keyList, key.Interface().(string))
	}
	sort.Strings(keyList)
	return
}
