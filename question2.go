package main

import (
	"container/heap"
	"fmt"
)

func main() {
	fmt.Println(rearrangeString("aab"))  // Example 1
	fmt.Println(rearrangeString("aaab")) // Example 2
}

type CharFrequency struct {
	char  rune
	count int
}

type MaxHeap []CharFrequency

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].count > h[j].count }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(CharFrequency))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func rearrangeString(s string) string {
	// Frequency map
	frequencyMap := make(map[rune]int)
	for _, ch := range s {
		frequencyMap[ch]++
	}

	// Create and populate the max heap
	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)
	for ch, count := range frequencyMap {
		heap.Push(maxHeap, CharFrequency{ch, count})
	}

	var result []rune
	var prev CharFrequency

	for maxHeap.Len() > 0 {
		current := heap.Pop(maxHeap).(CharFrequency)
		result = append(result, current.char)
		current.count--

		if prev.count > 0 {
			heap.Push(maxHeap, prev)
		}

		prev = current

		// If only one type of character is left and it's frequency is more than 1, it's not possible to rearrange
		if maxHeap.Len() == 1 && (*maxHeap)[0].count > 1 {
			return ""
		}
	}

	return string(result)
}
