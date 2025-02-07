package parser

import (
	// "bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"testing"
)

func setupLogging(logPath string) (err error) {
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	logOutputs := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logOutputs)
	log.SetFlags(0)

	return
}

func printXMLTree(decoder *xml.Decoder, indent string) error {
	var elemName string
	for {
		token, err := decoder.Token()
		if err != nil {
			return err // Return EOF or other errors
		}

		switch elem := token.(type) {
		case xml.StartElement:
			elemName = fmt.Sprintf("%s<%+v>", indent, elem)
			log.Println(elemName)
			// Recurse into nested elements
			if err := printXMLTree(decoder, indent+"  "); err != nil {
				if err.Error() == "EOF" {
					return nil
				}
				return err
			}
		case xml.EndElement:
			elemName = fmt.Sprintf("%s</%s>", indent[:len(indent)-2], elem.Name.Local)
			log.Println(elemName)
			return nil
		}
	}
}

func TestParseFeed(t *testing.T) {
	var testFiles = map[string]string{
		"https://archlinux.org/feeds/news/":                          path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "851066edb1ff2ed061430a9b89a3ab2657d9416f"),
		"https://www.pathofexile.com/news/rss":                       path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "3f2f2d4f33359839e533e70f5eb770fb1ba8d2b6"),
		"https://notrelated.xyz/rss":                                 path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "1dcd3a50b1a7b0f55b48e40b9a2babdbae932475"),
		"https://fractalsoftworks.com/feed/":                         path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "47f781c383cefb9f11cf37fc6d6ecebec92ac7d9"),
		"http://www.youtube.com/feeds/videos.xml?user=CaravanPalace": path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "a1c549e0bf1aee1f7c1c9858b5654352a62a3acf"),
	}
	for url, xmlFile := range testFiles {
		raw, err := os.ReadFile(xmlFile)
		if err != nil {
			t.Fatal(err)
		}
		feed, err := ParseFeed(raw, url)
		if err != nil {
			t.Fatal(err)
		}
		_ = feed
		// fmt.Println(feed.FeedItems[0].Content)
	}
}
