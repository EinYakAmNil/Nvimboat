package parser

import (
	// "bytes"
	"bytes"
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

func TestPrintFeed(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "851066edb1ff2ed061430a9b89a3ab2657d9416f")
	setupLogging("./test.log")
	raw, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatal(err)
	}
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	if err := printXMLTree(decoder, ""); err != nil && err.Error() != "EOF" {
		fmt.Println("Error parsing XML:", err)
	}
}

func TestParseFeed(t *testing.T) {
	var testFiles = map[string][]string{
		"https://archlinux.org/feeds/news/":                          {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "851066edb1ff2ed061430a9b89a3ab2657d9416f"), "https://archlinux.org/news/"},
		"https://www.pathofexile.com/news/rss":                       {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "3f2f2d4f33359839e533e70f5eb770fb1ba8d2b6"), "http://www.pathofexile.com"},
		"https://notrelated.xyz/rss":                                 {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "1dcd3a50b1a7b0f55b48e40b9a2babdbae932475"), "https://notrelated.xyz"},
		"https://fractalsoftworks.com/feed/":                         {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "47f781c383cefb9f11cf37fc6d6ecebec92ac7d9"), "https://fractalsoftworks.com"},
		"http://www.youtube.com/feeds/videos.xml?user=CaravanPalace": {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "a1c549e0bf1aee1f7c1c9858b5654352a62a3acf"), "https://www.youtube.com/channel/UCKH9HfYY_GEcyltl2mbD5lA"},
		"https://blog.lilydjwg.me/feed":                              {path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "b8afac7b449f2270a30ddd9b6a7c1fd5da4d75c4"), "https://blog.lilydjwg.me/"},
	}
	for url, testParams := range testFiles {
		raw, err := os.ReadFile(testParams[0])
		if err != nil {
			t.Fatal(err)
		}
		feed, err := ParseFeed(raw, url)
		if err != nil {
			t.Fatal(err)
		}
		if url != "https://www.pathofexile.com/news/rss" {
			if feed.FeedItems[0].Author == "" {
				t.Fatal("No author parsed", url)
			}
		}
		if feed.Url != testParams[1] {
			t.Fatalf(`feed.Url: expected "%s" got:"%s"`, testParams[1], feed.Url)
		}
	}
}
