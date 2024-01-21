package nvimboat

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func strings2bytes(stringSlice []string) [][]byte {
	byteSlices := [][]byte{}

	for _, s := range stringSlice {
		byteSlices = append(byteSlices, []byte(s))
	}

	return byteSlices
}

func parseFilterID(id string) (string, []string, []string, error) {
	var (
		query       string
		includeTags []string
		excludeTags []string
	)
	query, rawTags, _ := strings.Cut(id, ", ")
	_, query, _ = strings.Cut(query, "query: ")
	_, rawTags, _ = strings.Cut(rawTags, "tags: ")
	tags := strings.Split(rawTags, ", ")
	for _, t := range tags {
		if string(t[0]) == "!" {
			excludeTags = append(excludeTags, t[1:])
		} else {
			includeTags = append(includeTags, t)
		}
	}
	return query, includeTags, excludeTags, nil
}

func pageTypeString(p Page) string {
	fullName := fmt.Sprintf("%T", p)
	_, name, _ := strings.Cut(fullName, "nvimboat.")
	return name
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}
