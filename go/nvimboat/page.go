package nvimboat

func (ps *PageStack) Push(p *Page) {
	ps.Pages = append(ps.Pages, p)
}

func (ps *PageStack) Pop() {
	ps.Pages = ps.Pages[:len(ps.Pages)-1]
}
