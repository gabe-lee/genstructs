package genstructs

type BSNode[T any] struct {
	Parent   *BSNode[T]
	Children [2]*BSNode[T]
	Value    T
}

func (n *BSNode[T]) FirstAncestorOnSide(side NodeSide) *BSNode[T] {
	aSide := Left
	if side == Left {
		aSide = Right
	}
	parent := n.Parent
	node := n
	for parent != nil {
		pSide := node.Side()
		if pSide == aSide {
			return parent
		}
		node = parent
		parent = node.Parent
	}
	return parent
}

func (n *BSNode[T]) MaxBelow() *BSNode[T] {
	max := n
	for max.Children[Right] != nil {
		max = max.Children[Right]
	}
	return max
}

func (n *BSNode[T]) MinBelow() *BSNode[T] {
	min := n
	for min.Children[Left] != nil {
		min = min.Children[Left]
	}
	return min
}

func (n *BSNode[T]) Prev() *BSNode[T] {
	if n.Children[Left] == nil {
		return n.FirstAncestorOnSide(Left)
	}
	return n.Children[Left].MaxBelow()
}

func (n *BSNode[T]) Next() *BSNode[T] {
	if n.Children[Right] == nil {
		return n.FirstAncestorOnSide(Right)
	}
	return n.Children[Right].MinBelow()
}

func (n *BSNode[T]) Traverse(currentDepth uint64, op func(node *BSNode[T], depth uint64)) {
	if n.Children[Left] != nil {
		n.Children[Left].Traverse(currentDepth+1, op)
	}
	op(n, currentDepth)
	if n.Children[Right] != nil {
		n.Children[Right].Traverse(currentDepth+1, op)
	}
}

func (n *BSNode[T]) Side() NodeSide {
	if n.Parent.Children[Left] == n {
		return Left
	}
	return Right
}

func (n *BSNode[T]) OnLeft() bool {
	if n.Parent.Children[Left] == n {
		return true
	}
	return false
}

func (n *BSNode[T]) OnRight() bool {
	if n.Parent.Children[Right] == n {
		return true
	}
	return false
}

func (n *BSNode[T]) IsLeaf() bool {
	if n.Children[Right] == nil && n.Children[Left] == nil {
		return true
	}
	return false
}

func (n *BSNode[T]) HasLeft() bool {
	return n.Children[Left] != nil
}

func (n *BSNode[T]) HasRight() bool {
	return n.Children[Right] != nil
}

type BSTree[T any] struct {
	Root      *BSNode[T]
	NodeCount uint64
}

func (t *BSTree[T]) Traverse(op func(node *BSNode[T], depth uint64)) {
	if t.Root == nil {
		return
	}
	t.Root.Traverse(0, op)
}

func (t *BSTree[T]) Balance() {
	ordered := t.Flatten()
	t.Root = t.halfSplit(nil, ordered)
}

func (t *BSTree[T]) Add(value T, comparer CompareFunc[T]) (newNode *BSNode[T]) {
	if t.Root == nil {
		newNode = &BSNode[T]{
			Parent:   nil,
			Children: [2]*BSNode[T]{nil, nil},
			Value:    value,
		}
		t.Root = newNode
		t.NodeCount += 1
		return newNode
	}
	return t.AddFrom(t.Root, value, comparer)
}

func (t *BSTree[T]) AddFrom(parent *BSNode[T], value T, comparer CompareFunc[T]) (newNode *BSNode[T]) {
	if parent == nil {
		return nil
	}
	side := comparer(value, parent.Value)
	for parent.Children[side] != nil {
		parent = parent.Children[side]
		side = comparer(value, parent.Value)
	}
	newNode = &BSNode[T]{
		Parent:   parent,
		Children: [2]*BSNode[T]{nil, nil},
		Value:    value,
	}
	parent.Children[side] = newNode
	t.NodeCount += 1
	return newNode
}

func (t *BSTree[T]) Remove(node *BSNode[T], highSide NodeSide) {
	if node == nil {
		return
	}
	lowSide := Left
	if highSide == Left {
		lowSide = Right
	}
	parent := node.Parent
	side := node.Side()
	node.Children[lowSide].Parent = node.Children[highSide]
	node.Children[highSide].Parent = parent
	parent.Children[side] = node.Children[highSide]
	node = nil
	t.NodeCount -= 1
}

func (t *BSTree[T]) Find(value T, matcher MatchFunc[T]) (foundNode *BSNode[T]) {
	if t.Root == nil {
		return nil
	}
	foundNode = t.Root
	for {
		equal, side := matcher(value, foundNode.Value)
		if equal {
			return foundNode
		}
		foundNode = foundNode.Children[side]
		if foundNode == nil {
			return nil
		}
	}
}

func (t *BSTree[T]) Split(rootNode *BSNode[T], lowSide NodeSide, spliter SplitFunc[T]) (newRoot *BSNode[T]) {
	lVal, rVal := spliter(rootNode.Value)
	var lowVal, highVal T
	if lowSide == Left {
		lowVal = lVal
		highVal = rVal
	} else {
		lowVal = rVal
		highVal = lVal
	}
	rootNode.Value = highVal
	subNode := &BSNode[T]{
		Parent:   rootNode,
		Children: [2]*BSNode[T]{nil, nil},
		Value:    lowVal,
	}
	rootNode.Children[lowSide] = subNode
	subNode.Children[lowSide] = rootNode.Children[lowSide]
	subNode.Children[lowSide].Parent = subNode
	t.NodeCount += 1
	return rootNode
}

func (t *BSTree[T]) SplitAdd(value T, node *BSNode[T], lowSide NodeSide, spliter SplitFunc[T], comparer CompareFunc[T]) (newRoot *BSNode[T], newNode *BSNode[T]) {
	newRoot = t.Split(node, lowSide, spliter)
	return newRoot, t.AddFrom(newRoot, value, comparer)
}

func (t *BSTree[T]) Cull(node *BSNode[T], cullSide NodeSide) {
	orphan := node.Children[cullSide]
	orphanCount := uint64(0)
	orphan.Traverse(0, func(_ *BSNode[T], _ uint64) {
		orphanCount += 1
	})
	t.NodeCount -= orphanCount
	node.Children[cullSide] = nil
}

func (t *BSTree[T]) Flatten() []T {
	ordered := make([]T, 0, t.NodeCount)
	t.Traverse(func(node *BSNode[T], _ uint64) {
		ordered = append(ordered, node.Value)
	})
	return ordered
}

func (t *BSTree[T]) LeafImbalance() uint64 {
	lowest, highest := uint64(0), uint64((1<<64)-1)
	t.Traverse(func(node *BSNode[T], depth uint64) {
		if node.IsLeaf() {
			if depth > lowest {
				lowest = depth
			}
			if depth < highest {
				highest = depth
			}
		}
	})
	return lowest - highest
}

func (t *BSTree[T]) halfSplit(parent *BSNode[T], ordered []T) *BSNode[T] {
	if len(ordered) == 0 {
		return nil
	}
	mid := len(ordered) / 2
	node := &BSNode[T]{
		Parent:   parent,
		Children: [2]*BSNode[T]{nil, nil},
		Value:    ordered[mid],
	}
	lChild := t.halfSplit(node, ordered[:mid])
	node.Children[Left] = lChild
	mid += 1
	if mid >= len(ordered) {
		return node
	}
	rChild := t.halfSplit(node, ordered[mid:])
	node.Children[Right] = rChild
	return node
}
