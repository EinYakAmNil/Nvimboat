package nvimboat

import (
	"errors"
	"strconv"
)

func (f *Feed) MainPrefix() string {
	ratio := strconv.Itoa(f.UnreadCount) + "/" + strconv.Itoa(f.ArticleCount) + ")"
	if f.UnreadCount > 0 {

		return "N (" + ratio
	}
	return "  (" + ratio
}

func (f *Feed) PrefixCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Prefix())
	}
	return col
}

func (f *Feed) PubDateCol() ([]string, error) {
	var (
		col  []string
		err  error
		date string
	)
	for _, a := range f.Articles {
		date, err = unixToDate(a.PubDate)
		if err != nil {
			return nil, err
		}
		col = append(col, date)
	}
	return col, nil
}

func (f *Feed) AuthorCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Author)
	}
	return col
}

func (f *Feed) TitleCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Title)
	}
	return col
}

func (f *Feed) UrlCol() []string {
	var col []string

	for _, a := range f.Articles {
		col = append(col, a.Url)
	}
	return col
}

func (f *Feed) Render() ([][]string, error) {
	dates, err := f.PubDateCol()
	if err != nil {
		return nil, err
	}
	return [][]string{f.PrefixCol(), dates, f.AuthorCol(), f.TitleCol(), f.UrlCol()}, nil
}

func (f *Feed) ElementIdx(article Page) (int, error) {
	for i, a := range f.Articles {
		if a.Url == article.(*Article).Url {
			return i, nil
		}
	}
	return 0, errors.New("Couldn't find article in feed.")
}
