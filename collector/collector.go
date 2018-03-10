package collector

type Collector interface {
	Next() bool
	Collect() ([]Result, error)
}
