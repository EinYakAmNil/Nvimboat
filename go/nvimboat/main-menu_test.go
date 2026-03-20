package nvimboat

import (
	"reflect"
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
	mm := MainMenu{
		Feeds: []rssdb.QueryMainPageRow{
			{Title: "Abc"},
			{Title: "Abd"},
			{Title: "Bbc"},
			{Title: "abc"},
			{Title: "bbc"},
		},
	}
	for i, f := range mm.Feeds {
		idx, err := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Title: f.Title}})
		if err != nil {
			t.Fatal(err)
		}
		if mm.Feeds[i].Title != mm.Feeds[idx-len(FilterConfig)].Title {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx-len(FilterConfig)])
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
			"unread": int64(1),
			"tags":   []any{"!C", "B"},
		}, {
			"name":  "TestFilter3",
			"title": "abc",
			"tags":  []any{"!A", "C"},
		},
	}
	expectedFilters := []Filter{
		{
			Name: "TestFilter1",
			QueryFilterParams: rssdb.QueryFilterParams{
				Guid:            "%",
				Title:           "%",
				Author:          "%",
				Url:             "%",
				Content:         "%",
				UnreadStates:    []int{0, 1},
				ContentMimeType: "%",
			},
			FilterDescription: "filter: tags: A, !B",
			IncludeTags:       map[string]bool{"A": true},
			ExcludeTags:       map[string]bool{"B": true},
		},
		{
			Name: "TestFilter2",
			QueryFilterParams: rssdb.QueryFilterParams{
				Guid:            "%",
				Title:           "%",
				Author:          "%",
				Url:             "%",
				Content:         "%",
				UnreadStates:    []int{1},
				ContentMimeType: "%",
			},
			FilterDescription: "filter: unread: 1, tags: !C, B",
			IncludeTags:       map[string]bool{"B": true},
			ExcludeTags:       map[string]bool{"C": true},
		},
		{
			Name: "TestFilter3",
			QueryFilterParams: rssdb.QueryFilterParams{
				Guid:            "%",
				Title:           "abc",
				Author:          "%",
				Url:             "%",
				Content:         "%",
				UnreadStates:    []int{0, 1},
				ContentMimeType: "%",
			},
			FilterDescription: "filter: title: abc, tags: !A, C",
			IncludeTags:       map[string]bool{"C": true},
			ExcludeTags:       map[string]bool{"A": true},
		},
	}
	for i, raw := range rawFilters {
		parsedFilter, err := parseFilter(raw)
		if err != nil {
			t.Fatalf("nvimboat/TestParseFilter: %v\n", err)
		}
		if !reflect.DeepEqual(parsedFilter, expectedFilters[i]) {
			t.Fatalf(
				"nvimboat/TestParseFilter:\nExpected:\n%s\nGot:\n%s",
				prettyStruct(expectedFilters[i]),
				prettyStruct(parsedFilter),
			)
		}
	}
}
