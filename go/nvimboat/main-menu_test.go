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
		idx, err := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Title: f.Title}})
		if err != nil {
			t.Fatal(err)
		}
		if mm.Feeds[i] != mm.Feeds[idx] {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx])
		}
	}
}
