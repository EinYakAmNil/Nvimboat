package nvimboat

func (mm *MainMenu) Render() ([][]string, error) {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	for _, f := range mm.Filters {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Name)
		urlCol = append(urlCol, f.FilterID)
	}
	for _, f := range mm.Feeds {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}, nil
}

func (mm *MainMenu) ElementIdx(feed Page) (int, error) {
	switch feed.(type) {
	case *Filter:
		return 10, nil
	case *Feed:
		return 10, nil
	default:
		return 0, nil
	}
}
