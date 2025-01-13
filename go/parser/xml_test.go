package parser

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestParse(t *testing.T) {
	xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "ce3abe666d14c50974ef261a0db008082dbb561f")
	raw, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatal(err)
	}
	feed, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range feed.FeedItems {
		fmt.Println(i)
	}
}
