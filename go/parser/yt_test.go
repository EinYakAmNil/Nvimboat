package parser

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestParseYtFeed(t *testing.T) {
	xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "a1c549e0bf1aee1f7c1c9858b5654352a62a3acf")
	raw, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatal(err)
	}
	feed, articles, err := ParseYtFeed(raw, "")
	for _, a := range articles {
		fmt.Println(a.Unread)
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) == 0 {
		t.Fail()
	}
	if feed.Rssurl == "" {
		t.Fail()
	}
}
