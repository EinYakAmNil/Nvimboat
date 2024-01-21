package nvimboat

func (mm *MainMenu) Render() ([]string, error) {
	return []string{
		` |   (0/15)    | 3Dommaker                                 | query: unread = 1, tags: YouTube, !Music`,
		` | N (4/15)    | 3Dommaker                                 | query: unread = 1, tags: Music`,
		` | N (3/19)    | A Dummie                                  | https://www.youtube.com/feeds/videos.xml?channel_id=UCBzYAtjNpudOSCP5_8W1iAw`,
		` | N (5/15)    | A Jolly Wangcore                          | https://www.youtube.com/feeds/videos.xml?user=ajollywanker`,
		` |   (0/14)    | ADAM FRIENDED                             | https://www.youtube.com/feeds/videos.xml?channel_id=UCy6Q3wg-PgsgO2XtQxZpZEg`,
		` |   (0/38)    | Aintops - Topic                           | https://www.youtube.com/feeds/videos.xml?channel_id=UCbBlqVT59IUpsxNT9Dg1zIw`,
		` | N (11/32)   | Akie秋絵                                  | https://www.youtube.com/feeds/videos.xml?channel_id=UCs_JLrcQMNMHZgclnYwmAcQ`,
		` |   (0/18)    | Alice Magic - 有栖魔法                    | https://www.youtube.com/feeds/videos.xml?channel_id=UCz9zXlWxa0rbx9Lmk0pkBcg`,
	}, nil
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

func (f *TagsPage) Render() ([]string, error) {
	return []string{"tags page."}, nil
}

func (f *TagFeeds) Render() ([]string, error) {
	return []string{"tag feeds."}, nil
}
