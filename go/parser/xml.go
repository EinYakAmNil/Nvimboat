package parser

import (
	"encoding/xml"
	"fmt"
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
	Title   string `xml:"title"`
	Content string `xml:"description"`
}

type Feed struct {
	Rssurl string
	Url    string
	Title  string
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
	return
}
