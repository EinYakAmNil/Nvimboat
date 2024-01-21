package nvimboat

import (
	"fmt"
	"sort"
)

func (f *TagsPage) Render() ([][]string, error) {
	var (
		lines []string
		prefix string
	)
	for tag, count := range f.TagFeedCount {
		lines = append(lines, fmt.Sprintf("%s %s (%d)", prefix, tag, count))
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	return [][]string{lines}, nil
}

func (tf *TagFeeds) Render() ([][]string, error) {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	for _, f := range tf.Feeds {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}, nil
}
