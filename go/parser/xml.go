package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
)

type (
	Rss struct {
		Channel Channel `xml:"channel"`
	}
	Channel struct {
		Items []Item    `xml:"item"`
		Urls  []RssLink `xml:"link"`
		Title string    `xml:"title"`
	}
	RssLink struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	}
	Item struct {
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

func (raw Item) Item() (a rssdb.GetArticleRow, err error) {
	var (
		pubDate time.Time
	)
	const trims = "\n\t "
	timeFormats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
	}
	for _, format := range timeFormats {
		pubDate, err = time.Parse(format, raw.PubDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf(": %w\n"+
			"parser/Item", err,
		)
		return
	}
	a.Pubdate = pubDate.Unix()
	a.Author = strings.Trim(firstNonEmpty(raw.Author, raw.Creator),
		trims)
	a.Content = strings.Trim(firstNonEmpty(raw.Content, raw.Description, raw.Encoded), trims)
	a.Guid = raw.Guid
	a.Title = raw.Title
	a.Url = raw.Url
	return
}

func ParseDefaultFeed(raw []byte, url string) (
	feed *rssdb.InsertFeedParams,
	articles map[string]*rssdb.InsertArticleParams,
	err error,
) {
	var rss Rss
	feed = new(rssdb.InsertFeedParams)
	articles = make(map[string]*rssdb.InsertArticleParams)

	err = xml.Unmarshal(raw, &rss)
	if err != nil {
		err = fmt.Errorf("xml.Unmarshal: %w\n"+
			"parser/ParseDefaultFeed", err,
		)
		return
	}
	feed.Rssurl = url
findChannelUrl:
	for _, u := range rss.Channel.Urls {
		if u.Text != "" {
			feed.Url = u.Text
			break findChannelUrl
		}
	}
	for _, item := range rss.Channel.Items {
		a, itemErr := item.Item()
		if itemErr != nil {
			err = fmt.Errorf(": %w\n"+
				"parser/ParseDefaultFeed", itemErr,
			)
			return
		}
		articles[a.Guid] = &rssdb.InsertArticleParams{
			Author:  a.Author,
			Content: a.Content,
			Feedurl: url,
			Guid:    a.Guid,
			Pubdate: a.Pubdate,
			Title:   a.Title,
			Url:     a.Url,
		}
	}
	return
}
