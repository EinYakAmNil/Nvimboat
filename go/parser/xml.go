package parser

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName   xml.Name `xml:"channel"`
	Items     []Item   `xml:"item"`
	Links     []Links  `xml:"link"`
	Title     string   `xml:"title"`
	Desc      string   `xml:"description"`
	Generator string   `xml:"generator"`
}

type Links struct {
	XMLName xml.Name `xml:"link"`
	RssUrl  string   `xml:",chardata"`
	Url     string   `xml:"href,attr"`
}

type Item struct {
	Title   string   `xml:"title"`
	Content string   `xml:"description"`
	Guid    string   `xml:"guid"`
	Author  string   `xml:"author"`
	Url     string   `xml:"link"`
	Pubdate string   `xml:"pubDate"`
}

type Feed struct {
	Rssurl    string
	Url       string
	Title     string
	FeedItems []FeedItem
}

type FeedItem struct {
	rssdb.GetArticleRow
}

func Parse(raw []byte) (feed Feed, err error) {
	var rss Rss
	err = xml.Unmarshal(raw, &rss)
	if err != nil {
		err = fmt.Errorf("Parse: %w", err)
		return
	}
	feed.Title = rss.Channel.Title
	for _, link := range rss.Channel.Links {
		if len(link.RssUrl) > 0 {
			feed.Rssurl = link.RssUrl
		}
		if len(link.Url) > 0 {
			feed.Url = link.Url
		}
	}
	var (
		feedItem = FeedItem{rssdb.GetArticleRow{Feedurl: feed.Rssurl, Unread: 1}}
		pubDate  time.Time
	)
	for _, item := range rss.Channel.Items {
		feedItem.Author = item.Author
		feedItem.Guid = item.Guid
		pubDate, err = time.Parse(time.RFC1123, item.Pubdate)
		if err != nil {
			err = fmt.Errorf("Parse:\npubDate parsing: %w\ninput: %v\n", err, item)
			return
		}
		feedItem.Pubdate = pubDate.Unix()
		feedItem.Title = item.Title
		feedItem.Url = item.Url
		feed.FeedItems = append(feed.FeedItems, feedItem)
	}
	return
}
