package main

import (
	"testing"
)

func TestMinHeapConstructor(t *testing.T) {
	heap := Constructor()
	if heap == nil {
		t.Error("Expected heap to be initialized")
	}
	if heap.items == nil {
		t.Error("Expected heap items to be initialized")
	}
}

func TestMinHeapPush(t *testing.T) {
	heap := Constructor()
	heap.push(5)
	if heap.items[0] != 5 {
		t.Error("Expected 5 to be at the top of the heap")
	}
	heap.push(3)
	if heap.items[0] != 3 {
		t.Error("Expected 3 to be at the top of the heap")
	}
	heap.push(10)
	if heap.items[0] != 3 {
		t.Error("Expected 3 to be at the top of the heap")
	}
	heap.push(1)
	expectedHeap := []int{1, 3, 10, 5}
	for i, item := range heap.items {
		if item != expectedHeap[i] {
			t.Errorf("Expected %d to be at index %d", item, i)
		}
	}
}

func TestMinHeapPop(t *testing.T) {
	heap := Constructor()
	heap.push(5)
	heap.push(3)
	heap.push(10)
	heap.push(1)
	if val, err := heap.pop(); val != 1 && err != nil {
		t.Error("Expected 1 to be popped from the heap")
	}
	expectedHeap := []int{3, 5, 10}
	for i, item := range heap.items {
		if item != expectedHeap[i] {
			t.Errorf("Expected %d to be at index %d", item, i)
		}
	}
}

func TestMinHeapOneElementPop(t *testing.T) {
	heap := Constructor()
	heap.push(5)
	value, _ := heap.pop()
	if value != 5 {
		t.Error("Expected 5 to be popped from the heap")
	}
	if len(heap.items) != 0 {
		t.Error("Expected heap to be empty")
	}
}

func TestMinHeapNoElementPop(t *testing.T) {
	heap := Constructor()
	if _, err := heap.pop(); err == nil {
		t.Error("Expected an error when popping from an empty heap")
	}
	if heap.len() != 0 {
		t.Error("Expected heap to be empty")
	}
}

func TestMinHeapTop(t *testing.T) {
	heap := Constructor()
	heap.push(5)
	heap.push(3)
	heap.push(10)
	heap.push(1)
	if val, err := heap.top(); val != 1 && err != nil {
		t.Error("Expected 1 to be at the top of the heap")
	}
}