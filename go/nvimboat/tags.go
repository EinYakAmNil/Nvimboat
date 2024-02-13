package nvimboat

import (
	"errors"
	"fmt"
	"sort"
)

func (tp *TagsPage) Render(unreadOnly bool) ([][]string, error) {
	var (
		lines []string
		prefix string
	)
	for tag, count := range tp.TagFeedCount {
		lines = append(lines, fmt.Sprintf("%s %s (%d)", prefix, tag, count))
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	return [][]string{lines}, nil
}

func (tf *TagFeeds) Render(unreadOnly bool) ([][]string, error) {
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

func (tp *TagsPage) SubPageIdx(feed Page) (int, error) {
	return 0, nil
}

func (tf *TagFeeds) SubPageIdx(feed Page) (int, error) {
	for i, f := range tf.Feeds {
		if f.RssUrl == feed.(*Feed).RssUrl {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find feed in tag.")
}
