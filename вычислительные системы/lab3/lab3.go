package main

import (
	"log"
	"math"
	"time"

	"math/rand"
	// mpi "github.com/sbromberger/gompi"
)

const TREE_DEPTH = 10

type Node struct {
	Value       int
	Left, Right int
}

func generate(deep int) ([]Node, int) {
	num_nodes := int(math.Pow(2, float64(deep))) - 1
	tree := make([]Node, 0, num_nodes)
	/*
		{0 -1 -1}
		{0 1 2} {1 -1 -1} {2 -1 -1}
		{0 1 2} {1 3 4} {2 5 6} {3 -1 -1} {4 -1 -1} {5 -1 -1} {6 -1 -1}
	*/
	for i, j := 0, 1; i < num_nodes; i++ {
		if i >= int(math.Pow(2, float64(deep)-1))-1 {
			tree = append(tree, Node{
				Value: rand.Intn(10),
				Left:  -1,
				Right: -1,
			})
		} else {
			tree = append(tree, Node{
				Value: rand.Intn(10),
				Left:  j,
				Right: j + 1,
			})
			j += 2
		}
	}
	return tree, num_nodes
}

// In-order traversal
func walk_linear(tree []Node, callback func(Node)) {
	stack := make([]int, 0, int(math.Sqrt(float64(len(tree)+1))))
	index := 0
	for index != -1 || len(stack) != 0 {
		for index != -1 {
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

func main() {
	// generate tree
	tree, size := generate(TREE_DEPTH)
	log.Printf("tree size: %d\n", size)

	acc := 0
	start := time.Now()
	walk_linear(tree[0:3], func(n Node) { acc += n.Value })
	log.Printf("time: %d us, Total value: %d\n", time.Now().Sub(start).Microseconds(), acc)
}
