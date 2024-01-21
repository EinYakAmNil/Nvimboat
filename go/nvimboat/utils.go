package nvimboat

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func initDB(dbpath string) (*sql.DB, error) {
	var err error
	if fileExists(dbpath) {
		d, err := sql.Open("sqlite3", dbpath)
		if err != nil {
			return d, err
		}
		return d, nil
	}
	dbDir := path.Dir(dbpath)
	err = os.MkdirAll(dbDir, os.FileMode(0755))
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	_, err = d.Exec(createDB)
	if err != nil {
		return d, err
	}
	return d, err
}

func strings2bytes(stringSlice []string) [][]byte {
	byteSlices := [][]byte{}

	for _, s := range stringSlice {
		byteSlices = append(byteSlices, []byte(s))
	}

	return byteSlices
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

func subset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value] += 1
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

func parseFilterID(id string) (string, []string, []string, error) {
	var (
		query       string
		includeTags []string
		excludeTags []string
	)
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
	return query, includeTags, excludeTags, nil
}

func extracUrls(content string) []string {
	var links []string
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

func pageTypeString(p Page) string {
	fullName := fmt.Sprintf("%T", p)
	_, name, _ := strings.Cut(fullName, "nvimboat.")
	return name
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}

func anyToString(base []any) []string {
	var conv []string
	for _, a := range base {
		conv = append(conv, a.(string))
	}
	return conv
}

func filterTags(config []map[string]any, inTags, exTags []string) []any {
	feedurls := make(map[string]bool)
	var (
		urls []any
	)
	for _, feed := range config {
		if subset(inTags, anyToString(feed["tags"].([]any))) {
			feedurls[feed["rssurl"].(string)] = true
		}
	}
	for _, feed := range config {
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
	return urls
}

func articlesFilterQuery(query string, n int) string {
	const (
		prefix = `
		SELECT guid, title, author, url, feedurl, pubDate, content, unread
		FROM rss_item WHERE deleted = 0 AND feedurl in (?`
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

func (nb *Nvimboat) addColumn(col []string, separator string) error {
	currentLines, err := nb.plugin.Nvim.BufferLines(*nb.buffer, 0, -1, false)
	if err != nil {
		return err
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

	err = nb.SetLines(lines)
	if err != nil {
		return err
	}

	vcl, err := nb.virtColLens()
	if err != nil {
		return err
	}

	maxLineLen := slices.Max(vcl)

	for i, l := range lines {
		diff = maxLineLen - vcl[i]
		lines[i] = l + strings.Repeat(" ", diff)
	}

	err = nb.SetLines(lines)
	if err != nil {
		return err
	}

	return nil
}

func (nb *Nvimboat) virtColLens() ([]int, error) {
	evalResult := []int{}
	const virtcols = "map(range(1, line('$')), \"virtcol([v:val, '$'])\")"
	err := nb.plugin.Nvim.Eval(virtcols, &evalResult)
	if err != nil {
		return nil, err
	}

	return evalResult, err
}

func (nb *Nvimboat) trimTrail() error {
	currentLines, err := nb.plugin.Nvim.BufferLines(*nb.buffer, 0, -1, false)
	if err != nil {
		return err
	}

	var lines = []string{}
	for _, l := range currentLines {
		lines = append(lines, strings.TrimRight(string(l), " "))
	}

	err = nb.SetLines(lines)
	if err != nil {
		return err
	}

	return nil
}

func (nb *Nvimboat) parseFilters() ([]Filter, error) {
	var (
		filters []Filter
		f       Filter
	)
	for _, filter := range nb.ConfigFilters {
		if name, ok := filter["name"]; ok {
			f.Name = name.(string)
		} else {
			errMsg := fmt.Sprintf("Failed to parse: %v", filter)
			return filters, errors.New(errMsg)
		}
		if query, ok := filter["query"]; ok {
			f.Query = query.(string)
			f.FilterID = "query: " + query.(string) + ", tags: "
		} else {
			errMsg := fmt.Sprintf("Failed to parse: %v", filter)
			return filters, errors.New(errMsg)
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
		f = Filter{}
	}
	return filters, nil
}
