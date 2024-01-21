package nvimboat

import (
	"fmt"
	"sort"
)

func (f *TagsPage) Render() ([]string, error) {
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
	return lines, nil
}

func (f *TagFeeds) Render() ([]string, error) {
	return []string{"tag feeds."}, nil
}
