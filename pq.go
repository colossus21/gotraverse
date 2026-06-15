package gotraverse

// pqItem is a node queued for expansion at a given priority. seq records
// insertion order so that ties break deterministically (FIFO), making search
// output stable across runs.
type pqItem struct {
	node     string
	priority int
	seq      int
}

// priorityQueue is a min-heap of *pqItem implementing container/heap.Interface.
// Lower priority is popped first; equal priorities pop in insertion order.
type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return pq[i].seq < pq[j].seq
	}
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x any) {
	*pq = append(*pq, x.(*pqItem))
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[:n-1]
	return it
}
