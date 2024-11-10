package reload

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mmcdole/gofeed"
)

// Requests the URL if not found in cacheDir or if the modification time of the cache file is too old.
// The request will be cached in cacheDir.
// Indicates with the return value fromCache if cache was used.
func GetRss(url string, header http.Header, cacheTime time.Duration, cacheDir string) (feed *gofeed.Feed, fromCache bool, err error) {
	rssParser := gofeed.NewParser()
	cachePath := path.Join(cacheDir, hashUrl(url))
	fileStats, err := os.Stat(cachePath)
	if err != nil || time.Now().Sub(fileStats.ModTime()) > cacheTime {
		err = fmt.Errorf("%s is not cached", url)
		content, reqErr := requestUrl(url, header)
		if reqErr != nil {
			reqErr = fmt.Errorf("GetRss: %w", errors.Join(err, reqErr))
			return nil, false, reqErr
		}
		log.Println("requested", url)
		err = cacheUrl(url, cacheDir, content)
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, false, err
		}
		log.Println("cached", url)
		feed, err = rssParser.ParseString(string(content))
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, false, err
		}
		return feed, false, err
	} else {
		log.Printf("reading %s from cache\n", url)
		content, err := os.ReadFile(cachePath)
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, false, err
		}
		feed, err = rssParser.ParseString(string(content))
		if err != nil {
			err = fmt.Errorf("GetRss: %w", err)
			return nil, true, err
		}
		return feed, true, err
	}
}

func requestUrl(url string, header http.Header) (content []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		err = fmt.Errorf("requestUrl: failed to request %s, status code: %d\n", url, resp.StatusCode)
		return
	}
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("requestUrl: %w", err)
		return
	}
	return
}

func cacheUrl(url string, cacheDir string, content []byte) (err error) {
	fileName := hashUrl(url)
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		err = fmt.Errorf("cacheUrl: %w", err)
		return
	}
	err = os.WriteFile(path.Join(cacheDir, fileName), content, 0644)
	if err != nil {
		err = fmt.Errorf("cacheUrl: %w", err)
		return
	}
	return
}

func hashUrl(url string) (fileName string) {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	hashBytes := hasher.Sum(nil)
	fileName = hex.EncodeToString(hashBytes)
	fileName = path.Clean(fileName)
	return
}
