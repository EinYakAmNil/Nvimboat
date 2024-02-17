package nvimboat

import (
	"testing"
)

func TestArticleBack(t *testing.T) {
	filter := Filter{ArticleCount: 3, Articles: []*Article{
		{Title: "Article 1"},
		{Title: "Article 2"},
		{Title: "Article 3"},
	}}
	selectedA := filter.Articles[1]
	idx, err := filter.ChildIdx(selectedA)
	if err != nil {
		t.Fatal("No index.")
	}
	if idx != 1 {
		t.Fatal("Wrong index.")
	}
}
