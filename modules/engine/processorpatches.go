package engine

type nodePatchOpKind int

const (
	nodePatchOpSet nodePatchOpKind = iota
	nodePatchOpSetFlex
	nodePatchOpAddTag
	nodePatchOpClear
)

type nodePatchOp struct {
	kind     nodePatchOpKind
	node     *Node
	attr     Attribute
	values   AttributeValues
	flexInit []any
	tag      string
}

type NodePatchSet struct {
	ops []nodePatchOp
}

func (ps *NodePatchSet) Set(node *Node, attr Attribute, values ...AttributeValue) {
	ps.ops = append(ps.ops, nodePatchOp{
		kind:   nodePatchOpSet,
		node:   node,
		attr:   attr,
		values: append(AttributeValues(nil), values...),
	})
}

func (ps *NodePatchSet) SetFlex(node *Node, flexInit ...any) {
	ps.ops = append(ps.ops, nodePatchOp{
		kind:     nodePatchOpSetFlex,
		node:     node,
		flexInit: append([]any(nil), flexInit...),
	})
}

func (ps *NodePatchSet) AddTag(node *Node, tag string) {
	ps.ops = append(ps.ops, nodePatchOp{
		kind: nodePatchOpAddTag,
		node: node,
		tag:  tag,
	})
}

func (ps *NodePatchSet) Clear(node *Node, attr Attribute) {
	ps.ops = append(ps.ops, nodePatchOp{
		kind: nodePatchOpClear,
		node: node,
		attr: attr,
	})
}

func (ps *NodePatchSet) HasOperations() bool {
	return len(ps.ops) > 0
}

func (ps *NodePatchSet) Apply(ao *IndexedGraph) {
	for _, op := range ps.ops {
		switch op.kind {
		case nodePatchOpSet:
			op.node.Set(op.attr, op.values...)
		case nodePatchOpSetFlex:
			op.node.SetFlex(op.flexInit...)
		case nodePatchOpAddTag:
			op.node.Tag(op.tag)
		case nodePatchOpClear:
			op.node.Clear(op.attr)
		}
	}
}

type edgeDeltaOpKind int

const (
	edgeDeltaOpAdd edgeDeltaOpKind = iota
	edgeDeltaOpClear
	edgeDeltaOpSet
)

type edgeDeltaOp struct {
	kind       edgeDeltaOpKind
	from       *Node
	to         *Node
	edge       Edge
	edgeBitmap EdgeBitmap
	force      bool
	merge      bool
}

type EdgeDelta struct {
	ops []edgeDeltaOp
}

func (ed *EdgeDelta) Add(from, to *Node, edge Edge, force bool) {
	ed.ops = append(ed.ops, edgeDeltaOp{
		kind:  edgeDeltaOpAdd,
		from:  from,
		to:    to,
		edge:  edge,
		force: force,
	})
}

func (ed *EdgeDelta) Clear(from, to *Node, edge Edge) {
	ed.ops = append(ed.ops, edgeDeltaOp{
		kind: edgeDeltaOpClear,
		from: from,
		to:   to,
		edge: edge,
	})
}

func (ed *EdgeDelta) Set(from, to *Node, eb EdgeBitmap, merge bool) {
	ed.ops = append(ed.ops, edgeDeltaOp{
		kind:       edgeDeltaOpSet,
		from:       from,
		to:         to,
		edgeBitmap: eb,
		merge:      merge,
	})
}

func (ed *EdgeDelta) Apply(ao *IndexedGraph) {
	mutations := make([]nodeEdgeMutation, 0, len(ed.ops))
	for _, op := range ed.ops {
		switch op.kind {
		case edgeDeltaOpAdd:
			mutations = append(mutations, nodeEdgeMutation{
				From:  op.from,
				To:    op.to,
				Edge:  op.edge,
				Merge: true,
				Force: op.force,
			})
		case edgeDeltaOpClear:
			mutations = append(mutations, nodeEdgeMutation{
				From:  op.from,
				To:    op.to,
				Edge:  op.edge,
				Merge: true,
				Clear: true,
			})
		case edgeDeltaOpSet:
			mutations = append(mutations, nodeEdgeMutation{
				From:       op.from,
				To:         op.to,
				Edge:       NonExistingEdge,
				EdgeBitmap: op.edgeBitmap,
				Merge:      op.merge,
			})
		}
	}
	ao.applyIndexedEdgeMutations(ao.resolveEdgeMutations(mutations))
}
