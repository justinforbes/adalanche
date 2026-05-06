package engine

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type FrozenGraph struct {
	graph *IndexedGraph
}

func (g *IndexedGraph) Freeze() *FrozenGraph {
	return &FrozenGraph{graph: g}
}

func (fg *FrozenGraph) IndexedGraph() *IndexedGraph {
	return fg.graph
}

func (fg *FrozenGraph) Order() int {
	return fg.graph.Order()
}

func (fg *FrozenGraph) Root() *Node {
	return fg.graph.Root()
}

func (fg *FrozenGraph) Iterate(each func(o *Node) bool) {
	fg.graph.nodeMutex.RLock()
	nodes := fg.graph.nodes
	for _, n := range nodes {
		if !each(n) {
			fg.graph.nodeMutex.RUnlock()
			return
		}
	}
	fg.graph.nodeMutex.RUnlock()
}

func (fg *FrozenGraph) IterateParallel(each func(o *Node) bool, parallelFuncs int) {
	if parallelFuncs == 0 {
		parallelFuncs = runtime.NumCPU()
	}

	fg.graph.nodeMutex.RLock()
	nodes := fg.graph.nodes

	queue := make(chan *Node, parallelFuncs*2)
	var wg sync.WaitGroup
	var stop atomic.Bool

	for i := 0; i < parallelFuncs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for o := range queue {
				if !each(o) {
					stop.Store(true)
				}
			}
		}()
	}

	for i, o := range nodes {
		if i&0x3ff == 0 && stop.Load() {
			break
		}
		queue <- o
	}
	close(queue)
	fg.graph.nodeMutex.RUnlock()
	wg.Wait()
}

func (fg *FrozenGraph) Find(attribute Attribute, value AttributeValue) (*Node, bool) {
	return fg.graph.Find(attribute, value)
}

func (fg *FrozenGraph) FindMulti(attribute Attribute, value AttributeValue) (NodeSlice, bool) {
	return fg.graph.FindMulti(attribute, value)
}

func (fg *FrozenGraph) FindTwo(attribute Attribute, value AttributeValue, attribute2 Attribute, value2 AttributeValue) (*Node, bool) {
	return fg.graph.FindTwo(attribute, value, attribute2, value2)
}

func (fg *FrozenGraph) FindTwoMulti(attribute Attribute, value AttributeValue, attribute2 Attribute, value2 AttributeValue) (NodeSlice, bool) {
	return fg.graph.FindTwoMulti(attribute, value, attribute2, value2)
}

func (fg *FrozenGraph) DistinguishedParent(o *Node) (*Node, bool) {
	return fg.graph.DistinguishedParent(o)
}

func (fg *FrozenGraph) IterateEdges(node *Node, direction EdgeDirection, iter func(target *Node, ebm EdgeBitmap) bool) {
	if direction > In {
		return
	}

	index, ok := fg.graph.nodeLookup.Load(node)
	if !ok {
		return
	}

	for nodeIndex, edgeCombo := range fg.graph.edges[direction][index] {
		eb := fg.graph.edgeCombos[edgeCombo]
		target := fg.graph.nodes[nodeIndex]
		if !iter(target, eb) {
			return
		}
	}
}

func (fg *FrozenGraph) EdgeIteratorRecursive(node *Node, direction EdgeDirection, edgeMatch EdgeBitmap, excludemyself bool, goDeeperFunc func(source, target *Node, edge EdgeBitmap, depth int) bool) {
	seenObjects := make(map[*Node]struct{})
	if excludemyself {
		seenObjects[node] = struct{}{}
	}
	fg.edgeIteratorRecursive(node, direction, edgeMatch, goDeeperFunc, seenObjects, 1)
}

func (fg *FrozenGraph) edgeIteratorRecursive(node *Node, direction EdgeDirection, edgeMatch EdgeBitmap, goDeeperFunc func(source, target *Node, edge EdgeBitmap, depth int) bool, appliedTo map[*Node]struct{}, depth int) {
	fg.IterateEdges(node, direction, func(target *Node, edge EdgeBitmap) bool {
		if _, found := appliedTo[target]; found {
			return true
		}

		edgeMatches := edge.Intersect(edgeMatch)
		if edgeMatches.IsBlank() {
			return true
		}

		appliedTo[target] = struct{}{}
		if goDeeperFunc(node, target, edgeMatches, depth) {
			fg.edgeIteratorRecursive(target, direction, edgeMatch, goDeeperFunc, appliedTo, depth+1)
		}
		return true
	})
}
