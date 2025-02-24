// Project for concurrent skip list
// - [ X ] Make a basic skip list
// - [ ] Corrections
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
	for node != nil && len(node.NextOnLevel) != 0 {
		if node != sentinel {
			log.Printf("%d\n", node.Val)
		}
		node = node.NextOnLevel[level]
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
	dummy := &Node{Val: -10000, NextOnLevel: make(map[int]*Node)}
	prev := dummy

	// Construct a linked list first
	for _, num := range list {
		curr = &Node{Val: num, NextOnLevel: map[int]*Node{}, PrevOnLevel: map[int]*Node{}}
		curr.PrevOnLevel[0] = prev
		prev.NextOnLevel[0] = curr
		prev = curr
	}
	
	// Then, construct SkipList until maxLevel
	for level := 1; level < maxLevel; level++ {
		prev = dummy
		curr = dummy.NextOnLevel[level - 1]

		// Do we want to enforce setting some number on the level?
		for curr != nil {
			if rand.Float32() <= p {
				prev.NextOnLevel[level] = curr
				curr.PrevOnLevel[level] = prev
				prev = curr
			}
			curr = curr.NextOnLevel[level - 1]
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

	for level := sl.MaxLevel - 1; level >= 0; level-- {
		if insertByLevel[level] != nil && (level == 0 || rand.Float32() <= sl.P) {
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
	
	// rand.Seed(5)
	skipList := makeSkipList(nums, p, maxLevel)

	for i := 0; i < 10; i++ {
		num := rand.Intn(20)
		skipList.search(num)
		skipList.insert(num)
		skipList.search(num)
		skipList.delete(rand.Intn(20))
	}

	skipList.print()
}
