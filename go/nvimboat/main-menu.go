package nvimboat

func (mm *MainMenu) Render() ([][]string, error) {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	// for _, f := range mm.Filters {
	// 	prefixCol = append(prefixCol, f.MainPrefix())
	// 	titleCol = append(titleCol, f.Name)
	// 	urlCol = append(urlCol, f.FilterID)
	// }
	for _, f := range mm.Feeds {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}, nil
}

func (f *Article) Render() ([][]string, error) {
	return [][]string{{"article."}}, nil
}
