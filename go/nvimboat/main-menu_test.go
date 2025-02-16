package nvimboat

import (
	_ "fmt"
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func TestMainMenuChildIdx(t *testing.T) {
	mm := MainMenu{
		Feeds: []rssdb.MainPageFeed{
			{Feedurl: "Abc"},
			{Feedurl: "Abd"},
			{Feedurl: "Bbc"},
			{Feedurl: "abc"},
			{Feedurl: "bbc"},
		},
	}
	for i, f := range mm.Feeds {
		idx, err := mm.ChildIdx(&Feed{RssFeed: rssdb.RssFeed{Rssurl: f.Feedurl}})
		if err != nil {
			t.Fatal(err)
		}
		if mm.Feeds[i] != mm.Feeds[idx] {
			t.Fatal("expected:", mm.Feeds[i], "got:", mm.Feeds[idx])
		}
	}
}
