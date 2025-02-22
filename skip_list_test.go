package main

import (
	"math/rand"
	"testing"
)

func setupSkipList(t *testing.T) *SkipList {
	rand.Seed(42)
	nums := []int{1, 3, 4, 5, 6, 7, 8, 10}
	return makeSkipList(nums, 0.5, 4)
}

func TestSkipListConstructor(t *testing.T) {
	sl := setupSkipList(t)
	
	// Test basic properties
	if sl.MaxLevel != 4 {
		t.Errorf("Expected MaxLevel to be 4, got %d", sl.MaxLevel)
	}
	
	if sl.P != 0.5 {
		t.Errorf("Expected probability to be 0.5, got %f", sl.P)
	}
	
	// Test if sentinel node exists
	if sl.Sentinel == nil {
		t.Error("Sentinel node should not be nil")
	}
	
	// Test if base level contains all elements
	curr := sl.Sentinel.NextOnLevel[0]
	count := 0
	expected := []int{1, 3, 4, 5, 6, 7, 8, 10}
	
	for curr != nil {
		if curr.Val != expected[count] {
			t.Errorf("Expected value %d at position %d, got %d", expected[count], count, curr.Val)
		}
		count++
		curr = curr.NextOnLevel[0]
	}
	
	if count != len(expected) {
		t.Errorf("Expected %d elements in base level, got %d", len(expected), count)
	}
}

func TestSkipListSearchSuccessful(t *testing.T) {
	sl := setupSkipList(t)
	
	testCases := []struct {
		target int
		exists bool
	}{
		{1, true},
		{4, true},
		{10, true},
		{2, false},
		{9, false},
		{11, false},
	}
	
	for _, tc := range testCases {
		node, err := sl.search(tc.target)
		if tc.exists {
			if err != nil {
				t.Errorf("Expected to find %d, but got error: %v", tc.target, err)
			}
			if node.Val != tc.target {
				t.Errorf("Expected to find value %d, but got %d", tc.target, node.Val)
			}
		} else {
			if err == nil {
				t.Errorf("Expected not to find %d, but found it", tc.target)
			}
		}
	}
}

func TestSkipListSearchNotFound(t *testing.T) {
	sl := setupSkipList(t)
	
	nonExistentValues := []int{0, 2, 9, 11, 100}
	
	for _, val := range nonExistentValues {
		_, err := sl.search(val)
		if err == nil {
			t.Errorf("Expected error when searching for non-existent value %d", val)
		}
	}
}


func TestSkipListInsertSuccessful(t *testing.T) {
	sl := setupSkipList(t)
	
	testCases := []struct {
		insert   int
		expected bool
	}{
		{2, true},
		{9, true},
		{11, true},
	}
	
	for _, tc := range testCases {
		success, err := sl.insert(tc.insert)
		if err != nil {
			t.Errorf("Unexpected error inserting %d: %v", tc.insert, err)
		}
		if success != tc.expected {
			t.Errorf("Expected insert success %v for value %d, got %v", tc.expected, tc.insert, success)
		}
		
		// Verify the value can be found after insertion
		node, err := sl.search(tc.insert)
		if err != nil {
			t.Errorf("Could not find inserted value %d", tc.insert)
		}
		if node.Val != tc.insert {
			t.Errorf("Found incorrect value after insertion. Expected %d, got %d", tc.insert, node.Val)
		}
	}
}

func TestSkipListSearchTouchesCorrectNumElements(t *testing.T) {
}
