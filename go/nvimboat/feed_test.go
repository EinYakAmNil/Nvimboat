package nvimboat

import (
	"testing"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func TestFeedChildIdx(t *testing.T) {
	feed := Feed{
		Articles: []rssdb.GetFeedPageRow{
			{Title: "Article 1", Pubdate: 0},
			{Title: "Article 2", Pubdate: 5},
			{Title: "Article 3", Pubdate: 10},
			{Title: "Article 4", Pubdate: 11},
			{Title: "Article 5", Pubdate: 13},
			{Title: "Article 6", Pubdate: 15},
			{Title: "Article 7", Pubdate: 17},
		},
	}
	for i, a := range feed.Articles {
		idx := feed.ChildIdx(&Article{rssdb.GetArticleRow{Pubdate: a.Pubdate}})
		if feed.Articles[i].Title != feed.Articles[idx].Title {
			t.Fatal("expected:", a, "got:", feed.Articles[idx])
		}
	}
}
