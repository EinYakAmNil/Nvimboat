package nvimboat

import "strconv"

func (f *Feed) MainPrefix() string {
	ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
	if f.UnreadCount > 0 {

		return "N (" + ratio
	}

	return "  (" + ratio
}

func (f *Feed) Render() ([][]string, error) {
	return [][]string{{"feed."}}, nil
}
