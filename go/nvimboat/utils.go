package nvimboat

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
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

func pageTypeString(p Page) string {
	fullName := fmt.Sprintf("%T", p)
	_, name, _ := strings.Cut(fullName, "nvimboat.")
	return name
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}

func filterTags(config []map[string]any, inTags, exTags []string) []string {
	feedurls := make(map[string]bool)
	var urls []string
	for _, feed := range config {
		for _, tag := range feed["tag"].([]string) {
			if slices.Contains(inTags, tag) {
				feedurls[feed["rssurl"].(string)] = true
				continue
			}
		}
	}
	for _, feed := range config {
		for _, tag := range feed["tag"].([]string) {
			if slices.Contains(exTags, tag) {
				delete(feedurls, feed["rssurl"].(string))
				continue
			}
		}
	}
	return urls
}

func articlesFilterQuery(query string, n int) string {
	const (
		prefix = `
		SELECT guid, title, author, url, feedurl, pubDate, content, unread
		FROM rss_item WHERE deleted = 0 AND url in (?`
		suffix = `) ORDER BY pubDate DESC`
	)
	if n < 2 {
		return prefix + suffix
	}

	articleCount := strings.Repeat(", ?", n-1)
	return prefix + articleCount + suffix
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
