-- name: GetFeed :one
SELECT 
	rss_feed.*,
	COALESCE(feed_stats.article_count, 0) AS article_count,
	COALESCE(feed_stats.unread_count, 0) AS unread_count
FROM rss_feed
LEFT JOIN (
	SELECT feedurl,
		CAST(SUM(unread) AS INTEGER) AS unread_count,
		COUNT(*) AS article_count
	FROM rss_item
	WHERE deleted = 0
	GROUP BY feedurl
) feed_stats
ON rss_feed.rssurl = feed_stats.feedurl
WHERE rss_feed.rssurl = ?
ORDER BY rss_feed.title;

-- name: CreateFeed :one
INSERT INTO rss_feed (
	rssurl, url, title
	) VALUES (
	?, ?, ?
	)
RETURNING *;

-- name: GetArticle :one
SELECT guid, title, author, url, feedurl, pubDate, content, unread FROM rss_item
WHERE url = ? LIMIT 1;

-- name: SetArticlesRead :exec
UPDATE rss_item
SET unread = 0
WHERE url IN (sqlc.slice('url'));

-- name: SetArticlesUnread :exec
UPDATE rss_item
SET unread = 1
WHERE url IN (sqlc.slice('url'));

-- name: GetFeedPage :many
SELECT unread, pubDate, author, title, url FROM rss_item
WHERE feedurl = ?
AND deleted = 0
ORDER BY pubDate DESC;

-- name: AddArticle :exec
INSERT INTO rss_item (
	guid, title, author, url, feedurl, pubDate, content, unread, enclosure_url, flags, content_mime_type
	) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteArticles :exec
UPDATE rss_item
SET deleted = 1
WHERE url IN (sqlc.slice('url'));

-- name: DeleteFeedArticles :exec
UPDATE rss_item
SET deleted = 1
WHERE feedurl IN (sqlc.slice('feedurl'));

-- name: CleanupDeleted :exec
DELETE FROM rss_item
WHERE deleted = 1;

-- name: QueryMainPage :many
SELECT 
	rss_feed.title,
	rss_feed.rssurl,
	COALESCE(feed_articles.article_count, 0),
	COALESCE(feed_articles.unread_count, 0)
FROM rss_feed
LEFT JOIN (
	SELECT feedurl,
		CAST(SUM(unread) AS INTEGER) AS unread_count,
		COUNT(*) AS article_count
	FROM rss_item WHERE deleted = 0
	GROUP BY feedurl
) feed_articles
ON rss_feed.rssurl = feed_articles.feedurl
ORDER BY rss_feed.title;

-- name: QueryFilter :many
SELECT guid, title, author, url, feedurl, pubDate, content, unread FROM rss_item
WHERE guid LIKE ?
AND title LIKE ?
AND author LIKE ?
AND url LIKE ?
AND feedurl IN (sqlc.slice('feedurls'))
AND content LIKE ?
AND unread IN (sqlc.slice('unread_states'))
AND content_mime_type LIKE ?
AND deleted = 0
ORDER BY pubDate DESC;

-- name: SetFeedsRead :exec
UPDATE rss_item
SET unread = 0
WHERE feedurl IN (sqlc.slice('feedurl'));

-- name: SetFeedsUnread :exec
UPDATE rss_item
SET unread = 1
WHERE feedurl IN (sqlc.slice('feedurl'));
