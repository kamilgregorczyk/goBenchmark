package benchmarks

func merge(a []int, b []int) []int {
	var r = make([]int, len(a)+len(b))
	var i = 0
	var j = 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			r[i+j] = a[i]
			i++
		} else {
			r[i+j] = b[j]
			j++
		}
	}

	for i < len(a) {
		r[i+j] = a[i];
		i++
	}
	for j < len(b) {
		r[i+j] = b[j];
		j++
	}

	return r
}

func Mergesort(items []int) []int {
	if len(items) < 2 {
		return items
	}

	var middle = len(items) / 2

	aChan := make(chan []int)
	bChan := make(chan []int)
	go func() {
		aChan <- Mergesort(items[:middle])
		close(aChan)
	}()
	go func() {
		bChan <- Mergesort(items[middle:])
		close(bChan)
	}()
	return merge(<-aChan, <-bChan)
}
