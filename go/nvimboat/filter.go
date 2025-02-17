package nvimboat

type Filter struct {
	Name        string
	Query       string
	IncludeTags []string
	ExcludeTags []string
}
