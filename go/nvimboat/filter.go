package nvimboat

import "strconv"

func (f *Filter) Render() ([][]string, error) {
	return [][]string{{"filter."}}, nil
}

func (f *Filter) MainPrefix() string {
	ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
	if f.UnreadCount > 0 {

		return "N (" + ratio
	}

	return "  (" + ratio
}
