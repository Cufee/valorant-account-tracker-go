package types

type Account struct {
	ID       string
	Tag      string
	Name     string
	Username string
	LastRank Rank
}
type Rank struct {
	Tier  int
	Name  string
	Color string
	Icon  string
}
