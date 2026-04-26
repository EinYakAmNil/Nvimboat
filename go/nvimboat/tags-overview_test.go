package nvimboat

import (
	"fmt"
	"sort"
	"testing"
)

var (
	expectedLines = []string{
		"Tag A (3)",
		"Tag B (6)",
	}
)

func TestConstructPage(t *testing.T) {
	TagConfig = map[string][]string{
		"Tag A": {"a", "b", "c"},
		"Tag B": {"u", "v", "w", "x", "y", "z"},
	}
	var lines []string
	for tag, urls := range TagConfig {
		lines = append(lines, fmt.Sprintf(`%s (%d)`, tag, len(urls)))
	}
	sort.Strings(lines)
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Fatalf("Line %d should be '%s'. Got:\n%s", i+1, expectedLines[i], line)
		}
	}
}
