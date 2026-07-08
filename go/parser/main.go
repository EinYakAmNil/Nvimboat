package parser

import (
	"errors"
	"fmt"
	"log"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

func ParseFeed(raw []byte, url string) (
	feed *rssdb.InsertFeedParams,
	articles map[string]*rssdb.InsertArticleParams,
	err error,
) {
	parsers := []func([]byte, string) (
		*rssdb.InsertFeedParams,
		map[string]*rssdb.InsertArticleParams,
		error,
	){
		ParseDefaultFeed,
		ParseYtFeed,
	}
	for _, parser := range parsers {
		parsedFeed, parsedArticles, parseErr := parser(raw, url)
		if parseErr != nil {
			log.Println(parseErr)
			continue
		}
		if len(parsedArticles) > len(articles) {
			feed = parsedFeed
			articles = parsedArticles
		}
	}
	if len(articles) == 0 {
		err = fmt.Errorf(`Couldn't parse "%s" with available parsers`, url)
		err = errors.Join(err, errors.New("parser/ParseFeed"))
		return
	}
	return
}
