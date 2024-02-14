package nvimboat

type (
	Page interface {
		Render(bool) ([][]string, error)
		SubPageIdx(Page) (int, error)
	}
	PageStack struct {
		Pages []Page
	}
	MainMenu struct {
		Filters []*Filter
		Feeds   []*Feed
	}
	Filter struct {
		Name         string
		FilterID     string
		Query        string
		IncludeTags  []string
		ExcludeTags  []string
		UnreadCount  int
		ArticleCount int
		Articles     []*Article
	}
	Feed struct {
		Title        string
		RssUrl       string
		UnreadCount  int
		ArticleCount int
		Articles     []*Article
	}
	Article struct {
		Author  string
		Content string
		FeedUrl string
		Guid    string
		PubDate int
		Title   string
		Unread  int
		Url     string
	}
	TagsPage struct {
		Feeds        []map[string]any
		TagFeedCount map[string]int
	}
	TagFeeds struct {
		Tag   string
		Feeds []*Feed
	}
)

func (ps *PageStack) Push(p Page) {
	ps.Pages = append(ps.Pages, p)
}

func (ps *PageStack) Pop() {
	ps.Pages = ps.Pages[:len(ps.Pages)-1]
}

func (ps *PageStack) Top() Page {
	return ps.Pages[len(ps.Pages)-1]
}
