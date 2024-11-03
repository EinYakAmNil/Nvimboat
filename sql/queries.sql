-- name: GetFeed :one
SELECT * FROM rss_feed
WHERE rssurl = ? LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM rss_feed
ORDER BY rssurl;

-- name: CreateFeed :one
INSERT INTO rss_feed (
	rssurl, url, title
	) VALUES (
	?, ?, ?
	)
RETURNING *;

-- name: DeleteFeed :exec
DELETE FROM rss_feed
WHERE rssurl = ?;

-- name: DeleteFeedArticles :exec
DELETE FROM rss_item
WHERE feedurl = ?;
