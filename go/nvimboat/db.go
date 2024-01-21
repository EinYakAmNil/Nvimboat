package nvimboat

func (nb *Nvimboat) QueryArticle(url string) (Article, error) {
	var a Article
	return a, nil
}

func (nb *Nvimboat) QueryFeed(url string) (Feed, error) {
	var f Feed
	return f, nil
}

func (nb *Nvimboat) QueryFilter(query string, inTags, exTags []string) (Filter, error) {
	var f Filter
	return f, nil
}

func (nb *Nvimboat) QueryTags() (TagsPage, error) {
	var tp TagsPage
	tp.TagFeedCount = make(map[string]int)
	tp.Feeds = nb.ConfigFeeds
	for _, feed := range tp.Feeds {
		for _, tag := range feed["tags"].([]any) {
			tp.TagFeedCount[tag.(string)]++
		}
	}
	return tp, nil
}

func (nb *Nvimboat) QueryTagFeeds(tag string) (TagFeeds, error) {
	var tf TagFeeds
	return tf, nil
}
