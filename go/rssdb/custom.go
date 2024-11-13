package rssdb

import (
	"context"
	"fmt"
)

const mapArticles = `-- name: MapArticles :many
SELECT guid FROM rss_item
WHERE feedurl = ?
ORDER BY pubDate DESC
`

func (q *Queries) MapArticles(ctx context.Context, feedurl string) (map[string]bool, error) {
	fmt.Println(feedurl)
	rows, err := q.db.QueryContext(ctx, mapArticles, feedurl)
	if err != nil {
		return nil, err
	}
	items := make(map[string]bool)
	defer rows.Close()
	for rows.Next() {
		var guid string
		if err := rows.Scan(&guid); err != nil {
			return nil, err
		}
		items[guid] = true
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
