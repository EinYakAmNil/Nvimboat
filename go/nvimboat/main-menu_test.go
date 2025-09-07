package nvimboat

import (
	"fmt"
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func TestFilterQuery(t *testing.T) {
	filter := Filter{
		IncludeTags: map[string]bool{"tag1": true, "tag2": true},
		ExcludeTags: map[string]bool{"tag3": true, "tag4": true},
	}
	feeds := []Feed{{ // should be included in filter
		Tags:    map[string]bool{"tag1": true},
		RssFeed: rssdb.RssFeed{Rssurl: "https://example0.com/rss"},
	}, { // should be included in filter
		Tags:    map[string]bool{"tag2": true},
		RssFeed: rssdb.RssFeed{Rssurl: "https://example1.com/rss"},
	}, { // should not be included in filter
		Tags:    map[string]bool{"tag3": true},
		RssFeed: rssdb.RssFeed{Rssurl: "https://example2.com/rss"},
	}, { // should not be included in filter
		Tags:    map[string]bool{"tag1": true, "tag4": true},
		RssFeed: rssdb.RssFeed{Rssurl: "https://example3.com/rss"},
	}}
	urls := []string{}
filterFeed:
	for _, f := range feeds {
		for excTag := range filter.ExcludeTags {
			if f.Tags[excTag] == true {
				continue filterFeed
			}
		}
		for incTag := range filter.IncludeTags {
			if f.Tags[incTag] == true {
				urls = append(urls, f.Rssurl)
				continue filterFeed
			}
		}
	}
	filterQuery := `SELECT * FROM rss_item WHERE feedurl in ('%s') AND %s`
	expected := `SELECT * FROM rss_item WHERE feedurl in ('https://example0.com/rss', 'https://example1.com/rss') AND unread = 1`
	if filterQuery != expected {
// 		t.Fatalf(`
// expected:	%s
// got:		%s
// 		`, expected, filterQuery)
	}
}

func TestMainMenuChildIdx(t *testing.T) {
	dummyFilters := []*Filter{
		new(Filter),
		new(Filter),
		new(Filter),
	}
	mm := MainMenu{
		Filters: dummyFilters,
		Feeds: []MainPageFeed{
			{QueryMainPageRow: rssdb.QueryMainPageRow{Title: "Abc"}},
			{QueryMainPageRow: rssdb.QueryMainPageRow{Title: "Abd"}},
			{QueryMainPageRow: rssdb.QueryMainPageRow{Title: "Bbc"}},
			{QueryMainPageRow: rssdb.QueryMainPageRow{Title: "abc"}},
			{QueryMainPageRow: rssdb.QueryMainPageRow{Title: "bbc"}},
		},
	}
	for i, f := range mm.Feeds {
		idx, err := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Title: f.Title}})
		if err != nil {
			t.Fatal(err)
		}
		if mm.Feeds[i].Title != mm.Feeds[idx-len(mm.Filters)].Title {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx-len(mm.Filters)])
		}
	}
}

func TestParseFilter(t *testing.T) {
	rawFilters := []map[string]any{
		{
			"name": "TestFilter1",
			"tags": []any{"A", "!B"},
		}, {
			"name":   "TestFilter2",
			"unread": 1,
			"tags":   []any{"!C", "B"},
		}, {
			"name":  "TestFilter3",
			"title": "abc",
			"tags":  []any{"!A", "C"},
		},
	}
	parsedFilters, err := parseFilters(rawFilters)
	if err != nil {
		t.Fatalf("nvimboat/TestParseFilter: %v\n", err)
	}
	for _, i := range parsedFilters {
		fmt.Println(prettyStruct(i))
	}
}
