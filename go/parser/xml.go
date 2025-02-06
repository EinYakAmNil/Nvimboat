package parser

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type Rss struct {
	XMLName   xml.Name  `xml:""`
	Channel   Channel   `xml:"channel"`
	ChannelId ChannelId `xml:"channelId"`
}

type Channel struct {
	Items     []Item `xml:"item"`
	Entry     []Item `xml:"entry"`
	Links     []Link `xml:"link"`
	Title     string `xml:"title"`
	Desc      string `xml:"description"`
	Generator string `xml:"generator"`
}

type ChannelId struct {
	Items     []Item `xml:"item"`
	Entry     []Item `xml:"entry"`
	Links     []Link `xml:"link"`
	Title     string `xml:"title"`
	Desc      string `xml:"description"`
	Generator string `xml:"generator"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	RssUrl  string   `xml:",chardata"`
	Url     string   `xml:"href,attr"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
	Encoded     string `xml:"encoded"`
	Guid        string `xml:"guid"`
	Author      string `xml:"author"`
	Url         string `xml:"link"`
	Pubdate     string `xml:"pubDate"`
}

func ParseDefaultFeed(raw []byte) (feed Feed, err error) {
	var (
		feedItem = rssdb.GetArticleRow{Unread: 1}
		rss      Rss
		pubDate  time.Time
	)
	err = xml.Unmarshal(raw, &rss)
	if err != nil {
		err = fmt.Errorf("ParseDefaultFeed: %w", err)
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
	feedItem.Feedurl = feed.Rssurl
	for _, item := range rss.Channel.Items {
		feedItem.Author = item.Author
		feedItem.Guid = item.Guid
		if len(item.Content) > len(item.Encoded) {
			feedItem.Content = item.Content
		}
		if len(item.Description) > len(item.Content) {
			feedItem.Content = item.Description
		}
		if len(item.Encoded) > len(item.Description) {
			feedItem.Content = item.Encoded
		}
		pubDate, err = time.Parse(time.RFC1123, item.Pubdate)
		if err != nil {
			err = fmt.Errorf("ParseDefaultFeed:\npubDate parsing: %w\ninput: %v\n", err, item)
			return
		}
		feedItem.Pubdate = pubDate.Unix()
		feedItem.Title = item.Title
		feedItem.Url = item.Url
		feed.FeedItems = append(feed.FeedItems, feedItem)
	}
	return
}
