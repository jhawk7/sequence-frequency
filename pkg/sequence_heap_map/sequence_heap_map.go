package sequence_heap_map

import (
	"container/heap"
	"fmt"
)

func InitSequenceHeapMap() *SequenceHeapMap {
	//initialize sequence map and heap
	seqMap := make(map[string]*sequenceNode)
	seqHeap := &sequenceHeap{}
	heap.Init(seqHeap)
	return &SequenceHeapMap{
		seqMap:  seqMap,
		seqHeap: seqHeap,
	}
}

func (m *SequenceHeapMap) Upsert(sequence string) {
	//initialize/update the sequence in map heap - O(logn)
	if _, ok := m.seqMap[sequence]; !ok {
		seqNode := sequenceNode{
			value:    sequence,
			priority: 1,
		}
		m.seqMap[sequence] = &seqNode
		heap.Push(m.seqHeap, &seqNode)
	} else {
		seqNode := m.seqMap[sequence]
		m.seqHeap.update(seqNode, seqNode.value, (seqNode.priority + 1))
	}
}

func (m *SequenceHeapMap) Pop() (string, int, error) {
	//pop largest value off max heap and remove the entry in the map O(logn)
	if m.seqHeap.Len() == 0 {
		err := fmt.Errorf("error, sequence heap map is empty")
		return "", -1, err
	}
	seqNode := heap.Pop(m.seqHeap).(*sequenceNode)
	delete(m.seqMap, seqNode.value)
	return seqNode.value, seqNode.priority, nil
}

/*
container/heap interface implementation
https://pkg.go.dev/container/heap@go1.19.1
*/
func (seqHeap sequenceHeap) Len() int {
	return len(seqHeap)
}

func (seqHeap sequenceHeap) Less(i, j int) bool {
	//sequenceNode priority (frequency) at i should be greater than the priority at j in max heap
	return seqHeap[i].priority > seqHeap[j].priority
}

func (seqHeap sequenceHeap) Swap(i, j int) {
	seqHeap[i], seqHeap[j] = seqHeap[j], seqHeap[i]
	seqHeap[i].index = i
	seqHeap[j].index = j
}

func (seqHeap *sequenceHeap) Push(x interface{}) {
	//also modifies the heap's length, not just it's contents
	n := len(*seqHeap)
	seqNode := x.(*sequenceNode)
	seqNode.index = n
	*seqHeap = append(*seqHeap, seqNode)
}

func (seqHeap *sequenceHeap) Pop() interface{} {
	old := *seqHeap
	n := len(old)
	seqNode := old[n-1]
	old[n-1] = nil
	seqNode.index = -1
	*seqHeap = old[0 : n-1]
	return seqNode
}

func (seqHeap *sequenceHeap) update(seqNode *sequenceNode, value string, priority int) {
	seqNode.value = value
	seqNode.priority = priority
	heap.Fix(seqHeap, seqNode.index)
}
