package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
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
	Links     []Link
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
	Author  string
	Content string
	Guid    string
	Pubdate string
	Title   string
	Url     string
}

func (i *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "author", "creator":
				d.DecodeElement(&i.Author, &se)
			case "description", "content", "encoded":
				d.DecodeElement(&i.Content, &se)
			case "guid":
				d.DecodeElement(&i.Guid, &se)
			case "pubDate":
				d.DecodeElement(&i.Pubdate, &se)
			case "title":
				d.DecodeElement(&i.Title, &se)
			case "link":
				d.DecodeElement(&i.Url, &se)
			}
		case xml.EndElement:
			if se.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
}

func ParseDefaultFeed(raw []byte, url string) (feed Feed, err error) {
	var (
		feedItem     = rssdb.GetArticleRow{Unread: 1}
		rss          Rss
		pubDate      time.Time
		timeParseErr error
	)
	timeFormats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
	}
	err = xml.Unmarshal(raw, &rss)
	if err != nil {
		err = fmt.Errorf("ParseDefaultFeed: %w", err)
		return
	}
	feed.Title = strings.Trim(rss.Channel.Title, "\n\t ")
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
		feedItem.Author = strings.Trim(item.Author, "\n\t ")
		feedItem.Guid = strings.Trim(item.Guid, "\n\t ")
		feedItem.Content = strings.Trim(item.Content, "\n\t ")
		feedItem.Title = strings.Trim(item.Title, "\n\t ")
		feedItem.Url = strings.Trim(item.Url, "\n\t ")
	timeParseLoop:
		for _, layout := range timeFormats {
			pubDate, timeParseErr = time.Parse(layout, item.Pubdate)
			if timeParseErr == nil {
				break timeParseLoop
			}
		}
		if timeParseErr != nil {
			err = fmt.Errorf(`Could not parse "%s" in feed "%s" with available time formats`, item.Pubdate, url)
			return
		}
		feedItem.Pubdate = pubDate.Unix()
		feed.FeedItems = append(feed.FeedItems, feedItem)
	}
	return
}
