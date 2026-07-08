package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type (
	Rss2 struct {
		Channel Channel2 `xml:"channel"`
	}
	Channel2 struct {
		Items []ItemRaw `xml:"item"`
		Urls  []RssLink `xml:"link"`
		Title string    `xml:"title"`
	}
	RssLink struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	}
	ItemRaw struct {
		Author      string `xml:"author"`
		Creator     string `xml:"creator"`
		Content     string `xml:"content"`
		Description string `xml:"description"`
		Encoded     string `xml:"encoded"`
		Guid        string `xml:"guid"`
		PubDate     string `xml:"pubDate"`
		Title       string `xml:"title"`
		Url         string `xml:"link"`
	}
)

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func (rl RssLink) URL() string {
	if rl.Text != "" {
		return strings.TrimSpace(rl.Text)
	}
	return strings.TrimSpace(rl.Href)
}

func (raw ItemRaw) Item() Item {
	return Item{
		Author:  firstNonEmpty(raw.Author, raw.Creator),
		Content: firstNonEmpty(raw.Content, raw.Description, raw.Encoded),
	}
}

func ParseDefaultFeed2(raw []byte, url string) (feed Feed, err error) {
	var (
		rss          Rss2
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
		err = fmt.Errorf("xml.Unmarshal: %w\n"+
			"parser/ParseDefaultFeed2", err,
		)
		return
	}
findChannelUrl:
	for _, u := range rss.Channel.Urls {
		if u.Text != "" {
			feed.Url = u.Text
			break findChannelUrl
		}
	}
	for _, item := range rss.Channel.Items {
	timeParseLoop:
		for _, format := range timeFormats {
			pubDate, timeParseErr = time.Parse(format, item.PubDate)
			if timeParseErr == nil {
				break timeParseLoop
			}
		}
		if timeParseErr != nil {
			timeParseErr = fmt.Errorf(
				`Could not parse "%s" in feed "%s" with available time formats: %w`,
				item.PubDate,
				url,
				timeParseErr,
			)
			return
		}
		i := item.Item()
		feed.FeedItems = append(feed.FeedItems, rssdb.GetArticleRow{
			Author:  i.Author,
			Pubdate: pubDate.Unix(),
			Content: i.Content,
			Guid:    item.Guid,
			Title:   item.Title,
			Url:     item.Url,
		})
	}
	return
}
