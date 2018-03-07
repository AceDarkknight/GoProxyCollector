package collector

type Collector interface {
	Next() string
	Collect(url string) ([]Result, error)
}
