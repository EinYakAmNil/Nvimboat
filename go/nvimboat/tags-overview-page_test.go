package nvimboat

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

var (
	feedUrls = []string{
		"https://www.pathofexile.com/news/rss",
		"https://blog.lilydjwg.me/feed",
	}
	expectedFeeds = []rssdb.QueryTagFeedsRow{{
		Title:        "Path of Exile News",
		Feedurl:      feedUrls[0],
		UnreadCount:  34,
		ArticleCount: 37,
	}, {
		Title:        "依云's Blog",
		Feedurl:      feedUrls[1],
		UnreadCount:  12,
		ArticleCount: 12,
	}}
	tp = TagsOverviewPage{
		Tags: map[string][]string{
			"Tag A": {"a", "b", "c"},
			"Tag B": {"u", "v", "w", "x", "y", "z"},
		},
	}
	expectedLines = []string{
		"Tag A (3)",
		"Tag B (6)",
	}
)

func TestConstructPage(t *testing.T) {
	var lines []string
	for tag, urls := range tp.Tags {
		lines = append(lines, fmt.Sprintf(`%s (%d)`, tag, len(urls)))
	}
	sort.Strings(lines)
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Fatalf("Line %d should be '%s'. Got:\n%s", i+1, expectedLines[i], line)
		}
	}
}

func TestQueryTagFeeds(t *testing.T) {
	dbPath := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "tag-test.db")
	dbh, err := rssdb.ConnectDb(dbPath)
	if err != nil {
		t.Fatal()
	}
	feeds, err := dbh.Queries.QueryTagFeeds(dbh.Ctx, feedUrls)
	if err != nil {
		t.Fatal(err)
	}
	for i, feed := range feeds {
		if !reflect.DeepEqual(feed, expectedFeeds[i]) {
			t.Fatalf(
				"expected:\n%+v\ngot:\n%+v\n",
				feed,
				expectedFeeds[i],
			)
		}
	}
}
