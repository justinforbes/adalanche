package engine

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/lkarlslund/adalanche/modules/windowssecurity"
)

type frozenEdge struct {
	target NodeIndex
	edge   EdgeBitmap
}

type FrozenGraph struct {
	graph       *IndexedGraph
	root        *Node
	nodes       []*Node
	nodeIndexes map[*Node]NodeIndex
	edges       [2][][]frozenEdge
}

func (g *IndexedGraph) Freeze() *FrozenGraph {
	fg := &FrozenGraph{graph: g}

	g.nodeMutex.RLock()
	fg.root = g.root
	fg.nodes = append([]*Node(nil), g.nodes...)
	g.nodeMutex.RUnlock()

	fg.nodeIndexes = make(map[*Node]NodeIndex, len(fg.nodes))
	for i, node := range fg.nodes {
		fg.nodeIndexes[node] = NodeIndex(i)
	}

	g.edgeMutex.RLock()
	g.edgeComboMutex.RLock()
	for direction := range fg.edges {
		fg.edges[direction] = make([][]frozenEdge, len(fg.nodes))
		for from, toMap := range g.edges[direction] {
			if len(toMap) == 0 {
				continue
			}
			adjacency := make([]frozenEdge, 0, len(toMap))
			for target, edgeCombo := range toMap {
				adjacency = append(adjacency, frozenEdge{
					target: target,
					edge:   g.edgeCombos[edgeCombo],
				})
			}
			fg.edges[direction][from] = adjacency
		}
	}
	g.edgeComboMutex.RUnlock()
	g.edgeMutex.RUnlock()

	return fg
}

func (fg *FrozenGraph) IndexedGraph() *IndexedGraph {
	return fg.graph
}

func (fg *FrozenGraph) Order() int {
	return len(fg.nodes)
}

func (fg *FrozenGraph) Root() *Node {
	return fg.root
}

func (fg *FrozenGraph) Iterate(each func(o *Node) bool) {
	for _, n := range fg.nodes {
		if !each(n) {
			return
		}
	}
}

func (fg *FrozenGraph) IterateParallel(each func(o *Node) bool, parallelFuncs int) {
	if parallelFuncs == 0 {
		parallelFuncs = runtime.NumCPU()
	}

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

	for i, o := range fg.nodes {
		if i&0x3ff == 0 && stop.Load() {
			break
		}
		queue <- o
	}
	close(queue)
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

func (fg *FrozenGraph) FindAdjacentSID(s windowssecurity.SID, relativeTo *Node) (*Node, bool) {
	return fg.graph.FindAdjacentSID(s, relativeTo)
}

func (fg *FrozenGraph) IterateEdges(node *Node, direction EdgeDirection, iter func(target *Node, ebm EdgeBitmap) bool) {
	if direction > In {
		return
	}

	index, ok := fg.nodeIndexes[node]
	if !ok {
		return
	}

	for _, edge := range fg.edges[direction][index] {
		target := fg.nodes[edge.target]
		if !iter(target, edge.edge) {
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
