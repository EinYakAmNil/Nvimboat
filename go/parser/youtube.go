package parser

import (
	"encoding/xml"
	"errors"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type (
	YTFeed struct {
		Links   []YTLink  `xml:"link"`
		Title   string    `xml:"title"`
		Entries []YTEntry `xml:"entry"`
	}
	YTLink struct {
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	}
	YTEntry struct {
		Author  YTAuthor  `xml:"author"`
		Content YTContent `xml:"group"`
		Guid    string    `xml:"id"`
		Pubdate string    `xml:"published"`
		Title   string    `xml:"title"`
		Url     YTLink    `xml:"link"`
	}
	YTAuthor struct {
		Name string `xml:"name"`
		Uri  string `xml:"uri"`
	}
	YTContent struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
	}
)

func ParseYtFeed(xmlBytes []byte, url string) (
	feed *rssdb.InsertFeedParams,
	articles map[string]*rssdb.InsertArticleParams,
	err error,
) {
	var (
		pubDate time.Time
		ytFeed  YTFeed
	)
	feed = new(rssdb.InsertFeedParams)
	articles = make(map[string]*rssdb.InsertArticleParams)
	err = xml.Unmarshal(xmlBytes, &ytFeed)
	if err != nil {
		err = errors.Join(err, errors.New("parser/ParseYtFeed"))
		return
	}
	feed.Title = ytFeed.Title

	// YouTube shows http instead of https. This respects the users protocol
	feed.Rssurl = url

	for _, link := range ytFeed.Links {
		if link.Rel == "alternate" {
			feed.Url = link.Href
		}
	}
	for _, entry := range ytFeed.Entries {
		pubDate, err = time.Parse(time.RFC3339, entry.Pubdate)
		if err != nil {
			err = errors.Join(err, errors.New("parser/ParseYtFeed"))
			return
		}
		articles[entry.Guid] = &rssdb.InsertArticleParams{
			Pubdate: pubDate.Unix(),
			Feedurl: url,
			Author:  entry.Author.Name,
			Guid:    entry.Guid,
			Title:   entry.Title,
			Url:     entry.Url.Href,
			Content: entry.Content.Description,
		}
	}
	return
}
