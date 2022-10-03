package sequence_heap_map

type ISequenceHeapMap interface {
	Upsert(string)
	Pop() (string, int, error)
}

type SequenceHeapMap struct {
	//contains map of sequence strings to sequence node pointers and max heap of sequence nodes based on frequency (priority)
	seqMap  map[string]*sequenceNode
	seqHeap *sequenceHeap
}

type sequenceNode struct {
	//contains the sequence string and the count of occurrences
	value    string //the sequence (named 'value' to satisfy interface)
	priority int    //the frequency of the sequence (names 'priority' to satisfy interface)
	index    int    //the index will be automatically set by the container/heap pkg once the node is pushed
}

// SequenceHeap is a max heap of pointers to SequenceNodes
type sequenceHeap []*sequenceNode
