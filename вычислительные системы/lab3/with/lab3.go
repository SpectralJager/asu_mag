package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	mpi "github.com/sbromberger/gompi"
)

const TREE_DEPTH = 5

type Node struct {
	Value       int
	Left, Right int
}

type StartArgs struct {
	Tree       []Node
	StartIndex int
}

type DoneArgs struct {
	Result int
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

func main() {
	mpi.Start(false)

	depth := flag.Int("depth", TREE_DEPTH, "--depth: set maximum depth of tree")
	flag.Parse()

	// mpi workers communicator
	communicator := mpi.NewCommunicator(nil)
	rank := mpi.WorldRank()
	numProc := mpi.WorldSize()

	// parameters
	index := 0
	acc := 0
	tree := []Node{}

	if rank == 0 {
		tree, _ = generate(*depth)
		log.Printf("tree size: %d, Number of workers: %d\n", len(tree), numProc)
		switch numProc {
		case 3, 5, 9, 17, 33:
			firstDepth := int(math.Sqrt(float64(numProc)))
			childStart := calcNumNodes(firstDepth)
			for i := 1; i < numProc; i++ {
				args := StartArgs{
					Tree:       tree,
					StartIndex: childStart + i - 1,
				}
				buff := bytes.Buffer{}
				gob.NewEncoder(&buff).Encode(args)
				communicator.SendBytes(buff.Bytes(), i, i)
			}
			*depth = firstDepth
		case 1:
		default:
			panic("wrong number of processes")
		}
	} else {
		data, _ := communicator.RecvBytes(0, mpi.WorldRank())
		buff := bytes.Buffer{}
		buff.Write(data)
		var startArgs StartArgs
		gob.NewDecoder(&buff).Decode(&startArgs)
		tree = startArgs.Tree
		index = startArgs.StartIndex
	}

	start := time.Now()
	walk_linear(tree, index, *depth, func(n Node) {
		acc += n.Value
		// fmt.Println(n.Value)
		time.Sleep(time.Millisecond)
	})
	fmt.Printf("proc #%d done, result %d\n", mpi.WorldRank(), acc)

	if rank == 0 {
		for i := 1; i < numProc; i++ {
			data, _ := communicator.RecvBytes(i, 0)
			buff := bytes.Buffer{}
			buff.Reset()
			buff.Write(data)
			var doneArgs DoneArgs
			gob.NewDecoder(&buff).Decode(&doneArgs)
			acc += doneArgs.Result
		}
		log.Printf("time: %d us; Total value: %d\n", time.Now().Sub(start).Microseconds(), acc)
	} else {
		doneArgs := DoneArgs{
			Result: acc,
		}
		buff := bytes.Buffer{}
		gob.NewEncoder(&buff).Encode(doneArgs)
		communicator.SendBytes(buff.Bytes(), 0, 0)
	}

	mpi.Stop()
}
