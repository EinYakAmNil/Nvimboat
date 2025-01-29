package parser

import (
	"fmt"
	"log"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type Feed struct {
	Rssurl    string
	Url       string
	Title     string
	FeedItems []rssdb.GetArticleRow
}

func ParseFeed(raw []byte, url string) (feed Feed, err error) {
	var parsedFeed Feed

	parsers := []func([]byte) (Feed, error){
		ParseDefaultFeed,
		ParseYtFeed,
	}
	for _, parser := range parsers {
		parsedFeed, err = parser(raw)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(parsedFeed.FeedItems) > len(feed.FeedItems) {
			parsedFeed = feed
		}
	}
	if len(parsedFeed.FeedItems) == 0 {
		err = fmt.Errorf(`ParseFeed: couldn't parse "%s" with available parsers\n`, url)
		return
	}
	return
}
