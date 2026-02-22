package main

import "github.com/mattwiller/dblchk"

func main() {}

//export newFilter
func NewFilter() dblchk.Filter {
	return dblchk.NewFilter(0)
}

//export addToFilter
func AddToFilter(filter dblchk.Filter, element []byte) {
	filter.Add(element)
}

//export filterContains
func FilterContains(filter dblchk.Filter, element []byte) bool {
	return filter.MayContain(element)
}

//export resetFilter
func ResetFilter(filter dblchk.Filter) {
	filter.Reset()
}
