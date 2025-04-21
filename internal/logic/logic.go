package logic

func Filter(Slicecity []string) []string {
	FilterMap := make(map[string]struct{})
	for _, i := range Slicecity {
		FilterMap[i] = struct{}{}
	}

	FilterSlice := make([]string, 0, 10)
	for _, i := range Slicecity {
		if _, ok := FilterMap[i]; ok {
			FilterSlice = append(FilterSlice, i)
			delete(FilterMap, i)
		}
	}

	for i, j := 0, len(FilterSlice)-1; i < j; i, j = i+1, j-1 {
		FilterSlice[i], FilterSlice[j] = FilterSlice[j], FilterSlice[i]
	}

	return FilterSlice // мапа как фильтр
}
