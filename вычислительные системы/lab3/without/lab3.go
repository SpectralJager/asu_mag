package main

import (
	"flag"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"
	// mpi "github.com/sbromberger/gompi"
)

const TREE_DEPTH = 5

type Node struct {
	Value       int
	Left, Right int
}

func calcNumNodes(deep int) int {
	return int(math.Pow(2, float64(deep))) - 1
}

func calcDepth(numNodes int) int {
	return int(math.Sqrt(float64(numNodes + 1)))
}

func generate(deep int) ([]Node, int) {
	num_nodes := calcNumNodes(deep)
	tree := make([]Node, 0, num_nodes)
	/*
		{0 -1 -1}
		{0 1 2} {1 -1 -1} {2 -1 -1}
		{0 1 2} {1 3 4} {2 5 6} {3 -1 -1} {4 -1 -1} {5 -1 -1} {6 -1 -1}
	*/
	for i, j := 0, 1; i < num_nodes; i++ {
		if i >= int(math.Pow(2, float64(deep)-1))-1 {
			tree = append(tree, Node{
				// Value: rand.Intn(10),
				Value: i,
				Left:  -1,
				Right: -1,
			})
		} else {
			tree = append(tree, Node{
				// Value: rand.Intn(10),
				Value: i,
				Left:  j,
				Right: j + 1,
			})
			j += 2
		}
	}
	return tree, num_nodes
}

// In-order traversal
func walk_linear(tree []Node, start int, depth int, callback func(Node)) {
	stack := []int{}
	maxIndex := calcNumNodes(depth)
	index := start
	for (index <= maxIndex && index != -1) || len(stack) != 0 {
		for index != -1 && index < maxIndex {
			stack = append(stack, index)
			index = tree[index].Left
		}
		index = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node := tree[index]
		callback(node)
		index = node.Right
	}
}

func walk_parallel(tree []Node, numProc int, callback func(Node)) {
	switch numProc {
	case 3, 5, 9, 17, 33:
		wg := sync.WaitGroup{}
		firstDepth := int(math.Sqrt(float64(numProc)))
		wg.Add(1)
		go func() {
			// fmt.Println("0th worker")
			walk_linear(tree, 0, firstDepth, callback)
			wg.Done()
		}()
		childStart := calcNumNodes(firstDepth)
		for i := 0; i < numProc-1; i++ {
			wg.Add(1)
			go func() {
				// fmt.Printf("%dth worker\n", i+1)
				walk_linear(tree, childStart+i, calcDepth(len(tree)), callback)
				wg.Done()
			}()
		}
		wg.Wait()
	default:
		walk_linear(tree, 0, calcDepth(len(tree)), callback)
	}
}

func main() {
	depth := flag.Int("depth", TREE_DEPTH, "--depth: set maximum depth of tree")
	numWorkers := flag.Int("num", 1, "--num: set number of processes (workers)")
	flag.Parse()
	// generate tree
	tree, size := generate(*depth)
	log.Printf("tree size: %d, Number of workers: %d\n", size, *numWorkers)

	start := time.Now()
	var acc int64 = 0
	walk_parallel(tree, *numWorkers, func(n Node) {
		atomic.AddInt64(&acc, int64(n.Value))
		// fmt.Println(n.Value)
		time.Sleep(time.Millisecond)
	})
	log.Printf("time: %d us; Total value: %d\n", time.Now().Sub(start).Microseconds(), acc)
}
