package main

import "testing"

func TestQueryCraft(t *testing.T) {
	zero := articleReadUpdate(0)
	if zero != "" {
		t.Fatalf("Expected empty string. Got: %s", zero)
	}
	expectedOne := `UPDATE rss_item SET unread = ? WHERE url IN (?)`
	one := articleReadUpdate(1)
	if one != expectedOne {
		t.Fatalf("Expected: %s. Got: %s", expectedOne, one)
	}
	expectedTwo := `UPDATE rss_item SET unread = ? WHERE url IN (?, ?)`
	two := articleReadUpdate(2)
	if two != expectedTwo {
		t.Fatalf("Expected: %s. Got: %s", expectedTwo, two)
	}
	expectedMore := `UPDATE rss_item SET unread = ? WHERE url IN (?, ?, ?, ?, ?, ?)`
	more := articleReadUpdate(6)
	if more != expectedMore {
		t.Fatalf("Expected: %s. Got: %s", expectedMore, more)
	}
}
