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
	// xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "ce3abe666d14c50974ef261a0db008082dbb561f")
	// xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "47f781c383cefb9f11cf37fc6d6ecebec92ac7d9")
	xmlFile := path.Join(os.Getenv("HOME"), ".cache", "nvimboat-test", "a1c549e0bf1aee1f7c1c9858b5654352a62a3acf")
	raw, err := os.ReadFile(xmlFile)
	if err != nil {
		t.Fatal(err)
	}
	feed, err := ParseFeed(raw, "Caravan Palace")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(feed.FeedItems[0].Content)
	for _, i := range feed.FeedItems {
		fmt.Println(i.Guid)
	}
}
