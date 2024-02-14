package nvimboat

const (
	nvimboatState       = "return package.loaded.nvimboat."
	nvimboatEnable      = nvimboatState + "enable()"
	nvimboatDisable     = nvimboatState + "disable()"
	nvimboatConfig      = nvimboatState + "config"
	nvimboatFeeds       = nvimboatState + "feeds"
	nvimboatFilters     = nvimboatState + "filters"
	nvimboatPage        = nvimboatState + "page"
	nvimboatSetPageType = nvimboatState + "page.set(...)"
)

const (
	createDB = `
	CREATE TABLE google_replay (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	guid VARCHAR(64) NOT NULL,
	state INTEGER NOT NULL,
	ts INTEGER NOT NULL
	);
	CREATE TABLE metadata ( 
	db_schema_version_major INTEGER NOT NULL,
	db_schema_version_minor INTEGER NOT NULL
	);
	CREATE TABLE rss_feed ( 
	rssurl VARCHAR(1024) PRIMARY KEY NOT NULL,
	url VARCHAR(1024) NOT NULL,
	title VARCHAR(1024) NOT NULL ,
	lastmodified INTEGER(11) NOT NULL DEFAULT 0,
	is_rtl INTEGER(1) NOT NULL DEFAULT 0,
	etag VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE TABLE rss_item ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	guid VARCHAR(64) NOT NULL,
	title VARCHAR(1024) NOT NULL,
	author VARCHAR(1024) NOT NULL,
	url VARCHAR(1024) NOT NULL,
	feedurl VARCHAR(1024) NOT NULL,
	pubDate INTEGER NOT NULL,
	content VARCHAR(65535) NOT NULL,
	unread INTEGER(1) NOT NULL ,
	enclosure_url VARCHAR(1024),
	enclosure_type VARCHAR(1024),
	enqueued INTEGER(1) NOT NULL DEFAULT 0,
	flags VARCHAR(52),
	deleted INTEGER(1) NOT NULL DEFAULT 0,
	base VARCHAR(128) NOT NULL DEFAULT "",
	content_mime_type VARCHAR(255) NOT NULL DEFAULT "",
	enclosure_description VARCHAR(1024) NOT NULL DEFAULT "",
	enclosure_description_mime_type VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX idx_deleted ON rss_item(deleted);
	CREATE INDEX idx_feedurl ON rss_item(feedurl);
	CREATE INDEX idx_guid ON rss_item(guid);
	CREATE INDEX idx_lastmodified ON rss_feed(lastmodified);
	CREATE INDEX idx_rssurl ON rss_feed(rssurl);
	`
	feedsMinimalQuery = `
	SELECT rss_feed.title, c.* FROM rss_feed
	LEFT JOIN (
	SELECT a.feedurl, b.unreadCount, a.articleCount
	FROM (
	SELECT feedurl, COUNT(*) AS articleCount
	FROM rss_item
	GROUP BY feedurl
	) a
	LEFT JOIN (
	SELECT feedurl, sum(unread) AS unreadCount
	FROM rss_item
	GROUP BY feedurl
	) b
	ON a.feedurl = b.feedurl
	) c
	ON rss_feed.rssurl = c.feedurl
	ORDER BY rss_feed.title
	`
	articleQuery = `
	SELECT guid, title, author, feedurl, pubDate, content, unread
	FROM rss_item WHERE url = ? AND deleted = 0
	`
	feedQuery = `SELECT title FROM rss_feed WHERE rssurl = ?`
	feedArticlesQuery = `
	SELECT guid, title, author, url, feedurl, pubDate, content, unread FROM rss_item
	WHERE feedurl = ? AND deleted = 0
	ORDER BY pubDate DESC
	`
)
