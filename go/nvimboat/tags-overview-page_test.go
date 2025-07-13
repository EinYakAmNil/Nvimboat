package nvimboat

import (
	"fmt"
	"testing"
)

func TestConstructPage(t *testing.T) {
	tp := TagsOverviewPage{
		Tags: map[string][]string{
			"Tag A": {"a", "b", "c"},
			"Tag B": {"x", "y", "z"},
		},
	}
	var lines []string
	for tag, urls := range tp.Tags {
		lines = append(lines, fmt.Sprintf(`%s (%d)`, tag, len(urls)))
	}
	fmt.Println(lines)
}
