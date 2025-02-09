// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package rssdb

import (
	"database/sql"
)

type GoogleReplay struct {
	ID    int64
	Guid  string
	State int64
	Ts    int64
}

type MainPageFeed struct {
	Title        string
	Feedurl      string
	UnreadCount  int
	ArticleCount int
}

type Metadata struct {
	DbSchemaVersionMajor int64
	DbSchemaVersionMinor int64
}

type RssFeed struct {
	Rssurl       string
	Url          string
	Title        string
	Lastmodified int
	IsRtl        int
	Etag         string
}

type RssItem struct {
	ID                           int64
	Guid                         string
	Title                        string
	Author                       string
	Url                          string
	Feedurl                      string
	Pubdate                      int64
	Content                      string
	Unread                       int
	EnclosureUrl                 sql.NullString
	EnclosureType                sql.NullString
	Enqueued                     int
	Flags                        sql.NullString
	Deleted                      int
	Base                         string
	ContentMimeType              string
	EnclosureDescription         string
	EnclosureDescriptionMimeType string
}
