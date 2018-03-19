package cgminer

// Stats is a generic stats interface, returned by method Stats()
type Stats interface {
	S7() (*StatsS7, error)
	S9() (*StatsS9, error)
	D3() (*StatsD3, error)
	L3() (*StatsL3, error)
	T9() (*StatsT9, error)
}
