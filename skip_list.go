// Project for concurrent skip list
// - [ X ] Make a basic skip list
// - [ X ] Corrections
// - [ ] Add concurrency
// - [ ] Benchmark concurrent vs non-concurrent skip list

package main

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
)

type Node struct {
	Val int
	NextOnLevel map[int]*Node
	PrevOnLevel map[int]*Node
}

type SkipList struct {
	Sentinel *Node
	MaxLevel int
	P float32  // Acceptance probability
}

func logOnLevel(node *Node, level int) {
	sentinel := node
	log.Println("Level " + strconv.Itoa(level))
	for node != nil {
		if node != sentinel {
			log.Printf("%d\n", node.Val)
		}
		// Check if the node has a next node at this level before moving
		if next, exists := node.NextOnLevel[level]; exists {
			node = next
		} else {
			break
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func makeSkipList(list []int, p float32, maxLevel int) *SkipList {
	var curr *Node
	dummy := &Node{Val: -10000, NextOnLevel: make(map[int]*Node), PrevOnLevel: make(map[int]*Node)}
	prev := dummy

	// Construct a linked list at level 0 first
	for _, num := range list {
		curr = &Node{Val: num, NextOnLevel: map[int]*Node{}, PrevOnLevel: map[int]*Node{}}
		curr.PrevOnLevel[0] = prev
		prev.NextOnLevel[0] = curr
		prev = curr
	}
	
	for _, num := range list {
		// Find the node with this value
		node := dummy.NextOnLevel[0]
		for node != nil && node.Val != num {
			node = node.NextOnLevel[0]
		}
		
		if node == nil {
			continue // Should never happen
		}
		
		// Consider promotion for this node
		shouldPromote := true
		for level := 1; level < maxLevel && shouldPromote; level++ {
			shouldPromote = rand.Float32() <= p
			
			if shouldPromote {
				// Find the proper insertion point
				prev := dummy
				for {
					// Try to find a node at this level that comes before our node
					next := prev.NextOnLevel[level]
					
					// If no next or next is past our node in the list
					if next == nil || next.Val > node.Val {
						// Insert at this position
						node.PrevOnLevel[level] = prev
						node.NextOnLevel[level] = next
						prev.NextOnLevel[level] = node
						if next != nil {
							next.PrevOnLevel[level] = node
						}
						break
					}
					prev = next
				}
			}
		}
	}

	return &SkipList{Sentinel: dummy, MaxLevel: maxLevel, P: p}
}


func (sl *SkipList) print() {
	log.Println()
	for level := sl.MaxLevel - 1; level >= 0; level-- {
		logOnLevel(sl.Sentinel, level)
	}
	log.Println()
}


func (sl *SkipList) search(target int) (*Node, error) {
	level := sl.MaxLevel - 1
	curr := sl.Sentinel.NextOnLevel[level]
	
	start := sl.Sentinel.NextOnLevel[0]
	var end *Node

	// If level is empty, we go down until we find non-empty level
	for curr == nil && level >= 0 {
		curr = sl.Sentinel.NextOnLevel[level]
		level--
	}

	for start != end && level >= 0 {
		// log.Printf("Level: %d, Value: %d\n", level, curr.Val)
		if curr.Val > target {
			end = curr
			
			if curr.PrevOnLevel[level] == nil || curr.PrevOnLevel[level] == sl.Sentinel {
				level--
				continue
			}

			if abs(curr.PrevOnLevel[level].Val - target) < abs(curr.Val - target) { 
				curr = curr.PrevOnLevel[level]
			} else {
				level--
			}
		} else if curr.Val < target {
			start = curr
			
			if curr.NextOnLevel[level] == nil || curr.NextOnLevel[level] == end {
				level--
				continue
			}
			
			if abs(curr.NextOnLevel[level].Val - target) < abs(curr.Val - target) { 
				curr = curr.NextOnLevel[level]
			} else {
				level--
			}
		} else {
			return curr, nil
		}
	}

	return nil, errors.New("not found")
}

func insertNode(node *Node, nodeToInsert *Node, level int) {
	// Let's say we have 1 <-> 3, and node to insert is 2
	// 1) 2 -> 1.Next
	// 2) 1.Next.Prev -> 2
	// 3) 1.Next -> 2
	// 4) 2.Prev = 1
	nodeToInsert.NextOnLevel[level] = node.NextOnLevel[level]
	
	if node.NextOnLevel[level] != nil {
		node.NextOnLevel[level].PrevOnLevel[level] = nodeToInsert
	}
	
	node.NextOnLevel[level] = nodeToInsert
	nodeToInsert.PrevOnLevel[level] = node
}


func deleteNode(nodeToDelete *Node, level int) {
	if nodeToDelete.NextOnLevel[level] != nil {
		nodeToDelete.NextOnLevel[level].PrevOnLevel[level] = nodeToDelete.PrevOnLevel[level]
	}
	
	if nodeToDelete.PrevOnLevel[level] != nil {
		nodeToDelete.PrevOnLevel[level].NextOnLevel[level] = nodeToDelete.NextOnLevel[level]
	}
}

func (sl *SkipList) insert(value int) (bool, error) {
	insertByLevel := map[int]*Node{}

	level := sl.MaxLevel - 1
	curr := sl.Sentinel.NextOnLevel[level]

	start := sl.Sentinel.NextOnLevel[0]
	var end *Node

	for curr == nil && level >= 0{
		curr = sl.Sentinel.NextOnLevel[level]
		level--
	}

	for start != end && level >= 0 {
		log.Printf("Level: %d, Value: %d\n", level, curr.Val)
		if curr.Val > value {
			end = curr
			
			if curr.PrevOnLevel[level] == nil || curr.PrevOnLevel[level] == sl.Sentinel {
				insertByLevel[level] = sl.Sentinel
				level--
				continue
			}

			if abs(curr.PrevOnLevel[level].Val - value) < abs(curr.Val - value) { 
				curr = curr.PrevOnLevel[level]
			} else {
				insertByLevel[level] = curr.PrevOnLevel[level]
				level--
			}
		} else if curr.Val <= value {
			start = curr
			
			if curr.NextOnLevel[level] == nil || curr.NextOnLevel[level] == end {
				insertByLevel[level] = curr
				level--
				continue
			}
			
			if abs(curr.NextOnLevel[level].Val - value) < abs(curr.Val - value) { 
				curr = curr.NextOnLevel[level]
			} else {
				insertByLevel[level] = curr
				level--
			}
		}
	}
	
	nodeToInsert := &Node{Val: value, NextOnLevel: make(map[int]*Node), PrevOnLevel: make(map[int]*Node)}

	if insertByLevel[0] != nil {
		insertNode(insertByLevel[0], nodeToInsert, 0)
	}
	
	shouldPromote := true
	for level := 1; level < sl.MaxLevel; level++ {
		if !shouldPromote {
			break
		}
		
		shouldPromote = rand.Float32() <= sl.P
		
		if shouldPromote && insertByLevel[level] != nil {
			insertNode(insertByLevel[level], nodeToInsert, level)
		}
	}

	return true, nil
}

func (sl *SkipList) delete(value int) (bool, error) {
	deleteByLevel := map[int]*Node{}

	level := sl.MaxLevel - 1
	curr := sl.Sentinel.NextOnLevel[level]

	start := sl.Sentinel.NextOnLevel[0]
	var end *Node

	for curr == nil && level >= 0{
		curr = sl.Sentinel.NextOnLevel[level]
		level--
	}

	for start != end && level >= 0 {
		log.Printf("Level: %d, Value: %d\n", level, curr.Val)
		if curr.Val > value {
			end = curr
			
			if curr.PrevOnLevel[level] == nil || curr.PrevOnLevel[level] == sl.Sentinel {
				level--
				continue
			}

			if abs(curr.PrevOnLevel[level].Val - value) < abs(curr.Val - value) { 
				curr = curr.PrevOnLevel[level]
			} else {
				level--
			}
		} else if curr.Val <= value {
			if curr.Val == value {
				deleteByLevel[level] = curr
			}

			start = curr
			
			if curr.NextOnLevel[level] == nil || curr.NextOnLevel[level] == end {
				level--
				continue
			}
			
			if abs(curr.NextOnLevel[level].Val - value) < abs(curr.Val - value) { 
				curr = curr.NextOnLevel[level]
			} else {
				level--
			}
		}
	}

	deleted := false
	for level := sl.MaxLevel - 1; level >= 0; level-- {
		if deleteByLevel[level] != nil {
			deleted = true
			deleteNode(deleteByLevel[level], level)
		}
	}

	var err error
	if !deleted {
		err = errors.New("not found")
	}

	return deleted, err
}

func main() {
	nums := []int{1, 3, 4, 5, 6, 7, 8, 10}
	var p float32 = 0.5
	maxLevel := 4
	
	// Set a fixed seed for reproducible results
	rand.Seed(42)
	skipList := makeSkipList(nums, p, maxLevel)
	
	log.Println("=== Initial skip list ===")
	skipList.print()

	log.Println("=== Testing operations ===")
	for i := 0; i < 5; i++ {
		num := rand.Intn(20)
		log.Printf("Searching for %d", num)
		node, err := skipList.search(num)
		if err != nil {
			log.Printf("Not found: %d", num)
		} else {
			log.Printf("Found: %d", node.Val)
		}
		
		log.Printf("Inserting %d", num)
		skipList.insert(num)
		
		log.Printf("Searching for %d after insertion", num)
		node, err = skipList.search(num)
		if err != nil {
			log.Printf("ERROR: Not found after insertion: %d", num)
		} else {
			log.Printf("Found after insertion: %d", node.Val)
		}
		
		// Delete a specific value instead of random
		deleteVal := i * 2
		log.Printf("Deleting %d", deleteVal)
		deleted, _ := skipList.delete(deleteVal)
		if deleted {
			log.Printf("Deleted: %d", deleteVal)
		} else {
			log.Printf("Not found for deletion: %d", deleteVal)
		}
	}

	log.Println("=== Final skip list ===")
	skipList.print()
}
