package sequence_heap_map

import (
	"testing"
)

func Test_Upsert_Insert(t *testing.T) {
	testHeapMap := InitSequenceHeapMap()
	sequence := "test sequence string"

	testHeapMap.Upsert(sequence)
	if _, ok := testHeapMap.seqMap[sequence]; !ok {
		t.Error("error, expected sequence to be present in sequence map")
	}

	if testHeapMap.seqHeap.Len() != 1 {
		t.Error("error, expected size of sequence heap to be 1")
	}

	heap := testHeapMap.seqHeap
	if (*heap)[0].value != sequence {
		t.Error("error, unexpected value of seqence node within heap")
	}
}

func Test_Upsert_Update(t *testing.T) {
	testHeapMap := InitSequenceHeapMap()
	sequence := "test sequence string"

	testHeapMap.Upsert(sequence)
	testHeapMap.Upsert(sequence) //updating existing sequence

	if _, ok := testHeapMap.seqMap[sequence]; !ok {
		t.Error("error, expected sequence to be present in sequence map")
	}

	if testHeapMap.seqHeap.Len() != 1 {
		t.Error("error, expected size of sequence heap to be 1")
	}

	heap := testHeapMap.seqHeap
	if (*heap)[0].value != sequence {
		t.Error("error, unexpected value of seqence node within heap")
	}

	if (*heap)[0].priority != 2 {
		t.Error("error, unexpected value of sequence node priority/frequency within heap")
	}
}

func Test_Pop(t *testing.T) {
	testHeapMap := InitSequenceHeapMap()
	sequence := "test sequence string"

	testHeapMap.Upsert(sequence)
	testHeapMap.Upsert(sequence)
	poppedSequence, frequency, popErr := testHeapMap.Pop()

	if popErr != nil {
		t.Errorf("error, unexpected popErr %v", popErr)
	}

	if poppedSequence != sequence {
		t.Error("error, unexpected value of popped sequence")
	}

	if frequency != 2 {
		t.Error("error, unexpected value of sequence frequency returned")
	}

	if _, ok := testHeapMap.seqMap[sequence]; ok {
		t.Error("error, unexpected presence of deleted sequence in sequence map")
	}

	if testHeapMap.seqHeap.Len() != 0 {
		t.Error("error, expected size of sequence heap to be 0")
	}
}

func Test_Pop_PopErr(t *testing.T) {
	testHeapMap := InitSequenceHeapMap()
	sequence := "test sequence string"

	poppedSequence, frequency, popErr := testHeapMap.Pop()

	if popErr == nil {
		t.Error("error, expected popErr")
	}

	if poppedSequence != "" {
		t.Error("error, unexpected value of popped sequence")
	}

	if frequency != -1 {
		t.Error("error, unexpected value of sequence frequency returned")
	}

	if _, ok := testHeapMap.seqMap[sequence]; ok {
		t.Error("error, unexpected presence of deleted sequence in sequence map")
	}

	if testHeapMap.seqHeap.Len() != 0 {
		t.Error("error, expected size of sequence heap to be 0")
	}
}
