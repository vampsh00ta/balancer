package main

// IntHeap это минимальная куча целых чисел.
type BackendHeap []*Backend

// Len, Less, Swap для реализации интерфейса sort.Interface
func (h BackendHeap) Len() int           { return len(h) }
func (h BackendHeap) Less(i, j int) bool { return h[i].PingTime < h[j].PingTime }
func (h BackendHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *BackendHeap) Push(x interface{}) {
	// Push и Pop используют приемники указателей,
	// потому что они изменяют длину среза,
	// не только его содержимое.
	*h = append(*h, x.(*Backend))
}
func (h *BackendHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Этот пример вставляет несколько int в IntHeap,
// проверяет минимум,
// и удаляет их в порядке приоритета.
