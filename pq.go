package gotraverse

// pqItem is a node queued for expansion at a given priority. seq records
// insertion order so ties break deterministically (FIFO), keeping search output
// stable across runs.
type pqItem[N comparable] struct {
	node     N
	priority float64
	seq      int
}

// priorityQueue is a min-heap of *pqItem implementing container/heap.Interface.
// Lower priority is popped first; equal priorities pop in insertion order.
type priorityQueue[N comparable] []*pqItem[N]

func (pq priorityQueue[N]) Len() int { return len(pq) }

func (pq priorityQueue[N]) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return pq[i].seq < pq[j].seq
	}
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue[N]) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue[N]) Push(x any) {
	*pq = append(*pq, x.(*pqItem[N]))
}

func (pq *priorityQueue[N]) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return it
}
