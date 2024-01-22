package nvimboat

import "errors"

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
		for i, f := range mm.Filters {
			if feed.(*Filter).FilterID == f.FilterID {
				return i, nil
			}
		}
	case *Feed:
		for i, f := range mm.Feeds {
			if feed.(*Feed).RssUrl == f.RssUrl {
				return i + len(mm.Filters), nil
			}
		}
	default:
		return 0, nil
	}
	return 0, errors.New("Couldn't find feed/filter.")
}
