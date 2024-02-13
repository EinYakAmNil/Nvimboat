package nvimboat

func (a *Article) Prefix() string {
	if a.Unread == 1 {
		return "N"
	}
	return " "
}

func (a Article) Render(bool) ([][]string, error) {
	date, err := unixToDate(a.PubDate)
	if err != nil {
		return nil, err
	}
	lines := []string{
		"Feed: " + a.FeedUrl,
		"Title: " + a.Title,
		"Author: " + a.Author,
		"Date: " + date,
		"Link: " + a.Url,
		"== Article Begin ==",
	}
	content, err := renderHTML(a.Content)
	if err != nil {
		return nil, err
	}
	lines = append(lines, content...)
	lines = append(lines, "", "# Links")
	lines = append(lines, extracUrls(a.Content)...)

	return [][]string{lines}, nil
}

func (a *Article) SubPageIdx(Page) (int, error) {
	return 0, nil
}
