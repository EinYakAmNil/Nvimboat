package nvimboat

import (
	"fmt"
	"strings"
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func TestFilterQuery(t *testing.T) {
	filter := Filter{
		Query:       "unread = 1",
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
	filterQuery = fmt.Sprintf(filterQuery, strings.Join(urls, "', '"), filter.Query)
	expected := `SELECT * FROM rss_item WHERE feedurl in ('https://example0.com/rss', 'https://example1.com/rss') AND unread = 1`
	if filterQuery != expected {
		t.Fatalf(`
expected:	%s
got:		%s
		`, expected, filterQuery)
	}
}

func TestMainMenuChildIdx(t *testing.T) {
	dummyFilters := make(map[string]*Filter)
	dummyFilters["a"] = new(Filter)
	dummyFilters["b"] = new(Filter)
	mm := MainMenu{
		Filters: dummyFilters,
		Feeds: []MainPageFeed{
			{MainPageFeed: rssdb.MainPageFeed{Feedurl: "Abc"}},
			{MainPageFeed: rssdb.MainPageFeed{Feedurl: "Abd"}},
			{MainPageFeed: rssdb.MainPageFeed{Feedurl: "Bbc"}},
			{MainPageFeed: rssdb.MainPageFeed{Feedurl: "abc"}},
			{MainPageFeed: rssdb.MainPageFeed{Feedurl: "bbc"}},
		},
	}
	for i, f := range mm.Feeds {
		idx, err := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Rssurl: f.Feedurl}})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(idx)
		if mm.Feeds[i].Feedurl != mm.Feeds[idx-len(mm.Filters)].Feedurl {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx-len(mm.Filters)])
		}
	}
}
