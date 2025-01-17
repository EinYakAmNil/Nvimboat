package parser

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type (
	YTFeed struct {
		XMLName xml.Name  `xml:"feed"`
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

func ParseYtFeed(xmlBytes []byte) (feed Feed, err error) {
	var (
		feedItem = FeedItem{rssdb.GetArticleRow{Unread: 1}}
		pubDate  time.Time
		ytFeed   YTFeed
	)
	err = xml.Unmarshal(xmlBytes, &ytFeed)
	if err != nil {
		err = fmt.Errorf("ParseYtFeed: %v\n", err)
		return
	}
	feed.Title = ytFeed.Title
	for _, link := range ytFeed.Links {
		if link.Rel == "self" {
			feed.Rssurl = link.Href
		}
		if link.Rel == "alternate" {
			feed.Url = link.Href
		}
	}
	feedItem.Feedurl = feed.Rssurl
	for _, entry := range ytFeed.Entries {
		pubDate, err = time.Parse(time.RFC3339, entry.Pubdate)
		if err != nil {
			err = fmt.Errorf("ParseYtFeed: bad pubDate\n%v: %v\n", err, entry.Pubdate)
		}
		feedItem.Pubdate = pubDate.Unix()
		feedItem.Author = entry.Author.Name
		feedItem.Guid = entry.Guid
		feedItem.Title = entry.Title
		feedItem.Url = entry.Url.Href
		feedItem.Content = entry.Content.Description
		feed.FeedItems = append(feed.FeedItems, feedItem)
	}
	return
}
