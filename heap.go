package main

import (
	"errors"
)

type MinHeap struct {
    items []int
}

func Constructor() *MinHeap {
	return &MinHeap{items: []int{}}
}

func (heap *MinHeap) push(item int) {
	heap.items = append(heap.items, item)
	heap.heapifyUp(len(heap.items) - 1)
}

func (heap *MinHeap) heapifyUp(index int) {
	parentIndex := (index - 1) / 2
	if index > 0 && heap.items[index] < heap.items[parentIndex] {
		heap.items[index], heap.items[parentIndex] = heap.items[parentIndex], heap.items[index]
		heap.heapifyUp(parentIndex)
	}
}

func (heap *MinHeap) len() int {
	return len(heap.items)
}

func (heap *MinHeap) top() (int, error) {
	if heap.len() == 0 {
		return 0, errors.New("Empty heap")
	}
	return heap.items[0], nil
}

func (heap *MinHeap) pop() (int, error)  {
	if heap.len() == 0 {
		return 0, errors.New("No value to pop")
	}
	if heap.len() == 1 {
		last := heap.items[len(heap.items) - 1]
		heap.items = heap.items[:len(heap.items) - 1]
		return last, nil
	}
	
	rootVal := heap.items[0]
	last := heap.items[len(heap.items) - 1]
	heap.items = heap.items[:len(heap.items) - 1]
	heap.items[0] = last
	
	heap.heapifyDown(0)
	
	return rootVal, nil
}

func (heap *MinHeap) heapifyDown(index int) {
	smallest := index
	leftChild := index * 2 + 1
	rightChild := index * 2 + 2
	
	if leftChild < heap.len() && heap.items[leftChild] < heap.items[smallest] {
		smallest = leftChild
	}
	
	if rightChild < heap.len() && heap.items[rightChild] < heap.items[smallest] {
		smallest = rightChild
	}
	
	if smallest != index {
		heap.items[smallest], heap.items[index] = heap.items[index], heap.items[smallest]
		heap.heapifyDown(smallest)
	}
}
