package extension

import "github.com/laz/dubbo-go/filter"

var (
	filters = make(map[string]func() filter.Filter)
)

// SetFilter sets the filter extension with @name
// For example: hystrix/metrics/token/tracing/limit/...
func SetFilter(name string, v func() filter.Filter) {
	filters[name] = v
}

// GetFilter finds the filter extension with @name
func GetFilter(name string) filter.Filter {
	if filters[name] == nil {
		panic("filter for " + name + " is not existing, make sure you have imported the package.")
	}
	return filters[name]()
}
