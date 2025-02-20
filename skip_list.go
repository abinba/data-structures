// Project for concurrent skip list
// - [ ] Make a basic skip list
// - [ ] Add concurrency
// - [ ] Benchmark concurrent vs non-concurrent skip list

package main

import (
	"errors"
	"fmt"
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
	log.Println("Level " + strconv.Itoa(level))
	for node != nil && len(node.NextOnLevel) != 0 {
		if node.Val != -1 {
			log.Printf("%d\n", node.Val)
		}
		node = node.NextOnLevel[level]
	}
}

func makeSkipList(list []int, p float32, maxLevel int) *SkipList {
	var curr *Node
	dummy := &Node{Val: -1, NextOnLevel: make(map[int]*Node)}
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
		
		first := true

		for curr != nil {
			if rand.Float32() <= p || first {
				prev.NextOnLevel[level] = curr
				curr.PrevOnLevel[level] = prev
				prev = curr
				first = false
			}
			curr = curr.NextOnLevel[level - 1]
		}
	}

	logOnLevel(dummy, 0)
	logOnLevel(dummy, 1)
	logOnLevel(dummy, 2)
	logOnLevel(dummy, 3)
	log.Println()

	return &SkipList{Sentinel: dummy, MaxLevel: maxLevel, P: p}
}


func (sl *SkipList) search(target int) (*Node, error) {
	level := sl.MaxLevel - 1
	curr := sl.Sentinel.NextOnLevel[level]
	
	start := curr
	var end *Node

	for start != end && level >= 0 {
		log.Printf("Level: %d, Value: %d\n", level, curr.Val)
		if curr.Val > target {
			end = curr
			
			if level > 0{
				level--
			}
			
			curr = curr.PrevOnLevel[level]
		} else if curr.Val < target {
			start = curr
			
			if curr.NextOnLevel[level] == nil || curr.NextOnLevel[level] == end {
				level--
				continue
			}

			curr = curr.NextOnLevel[level]
		} else {
			return curr, nil
		}
	}

	return nil, errors.New("not found")
}

// func (sl *SkipList) insert(value int) (bool, error) {

// }

// func (sl *SkipList) delete(value int) (bool, error) {

// }

func main() {
	nums := []int{1, 3, 4, 5, 6, 7, 8, 10}
	var p float32 = 0.5
	maxLevel := 4

	skipList := makeSkipList(nums, p, maxLevel)
	fmt.Println(skipList.search(0))
}
