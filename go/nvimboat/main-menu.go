package nvimboat

func (mm *MainMenu) Render() ([][]string, error) {
	var (
		prefixCol []string
		titleCol  []string
		urlCol    []string
	)
	for _, f := range mm.Feeds {
		prefixCol = append(prefixCol, f.MainPrefix())
		titleCol = append(titleCol, f.Title)
		urlCol = append(urlCol, f.RssUrl)
	}
	return [][]string{prefixCol, titleCol, urlCol}, nil
}

func (f *Filter) Render() ([]string, error) {
	return []string{"filter."}, nil
}

func (f *Feed) Render() ([]string, error) {
	return []string{"feed."}, nil
}

func (f *Article) Render() ([]string, error) {
	return []string{"article."}, nil
}
