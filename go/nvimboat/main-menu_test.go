package nvimboat

import (
	_ "fmt"
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func TestMainMenuChildIdx(t *testing.T) {
	mm := MainMenu{
		Feeds: []rssdb.MainPageFeed{
			{Title: "Abc"},
			{Title: "Abd"},
			{Title: "Bbc"},
			{Title: "abc"},
			{Title: "bbc"},
		},
	}
	for i, f := range mm.Feeds {
		idx := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Title: f.Title}})
		if mm.Feeds[i] != mm.Feeds[idx] {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx])
		}
	}
}
