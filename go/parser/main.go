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
	var (
		parsedFeed Feed
		parseErr error
	)

	parsers := []func([]byte) (Feed, error){
		ParseDefaultFeed,
		ParseYtFeed,
	}
	for _, parser := range parsers {
		parsedFeed, parseErr = parser(raw)
		if parseErr != nil {
			log.Println(parseErr)
			continue
		}
		if len(parsedFeed.FeedItems) > len(feed.FeedItems) {
			feed = parsedFeed
		}
	}
	if len(parsedFeed.FeedItems) == 0 {
		err = fmt.Errorf(`ParseFeed: couldn't parse "%s" with available parsers\n`, url)
		return
	}
	return
}
