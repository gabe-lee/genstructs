package genstructs

type RBColor uint8

const (
	Red   RBColor = 0
	Black RBColor = 1
)

type RBNode[T any] struct {
	Color    RBColor
	Parent   *RBNode[T]
	Children [2]*RBNode[T]
	Value    T
}

func (n *RBNode[T]) Side() NodeSide {
	if n.OnLeft() {
		return Left
	}
	return Right
}

func (n *RBNode[T]) OnLeft() bool {
	return n.Parent.Children[Left] == n
}

func (n *RBNode[T]) OnRight() bool {
	return n.Parent.Children[Right] == n
}

func (n *RBNode[T]) GrandParent() *RBNode[T] {
	return n.Parent.Parent
}

func (n *RBNode[T]) Uncle() *RBNode[T] {
	return n.Parent.Sibling()
}

func (n *RBNode[T]) Sibling() *RBNode[T] {
	return n.Parent.Children[1-n.Side()]
}

func (n *RBNode[T]) IsRed() bool {
	return n.Color == Red
}

func (n *RBNode[T]) IsBlack() bool {
	return n.Color == Black
}

func (n *RBNode[T]) SwapWith(other *RBNode[T]) {
	holdN := RBNode[T]{
		Color:    other.Color,
		Parent:   other.Parent,
		Children: [2]*RBNode[T]{other.Children[Left], other.Children[Right]},
	}
	holdO := RBNode[T]{
		Color:    other.Color,
		Parent:   other.Parent,
		Children: [2]*RBNode[T]{other.Children[Left], other.Children[Right]},
	}
	n.Children[Left].Parent = other
	n.Children[Right].Parent = other
	n.Parent.Children[n.Side()] = other
	other.Children[Left].Parent = n
	other.Children[Right].Parent = n
	other.Parent.Children[other.Side()] = n
	other.Color = holdN.Color
	other.Parent = holdN.Parent
	other.Children = [2]*RBNode[T]{holdN.Children[Left], holdN.Children[Right]}
	n.Color = holdO.Color
	n.Parent = holdO.Parent
	n.Children = [2]*RBNode[T]{holdO.Children[Left], holdO.Children[Right]}
}

type RBTree[T any] struct {
	Root *RBNode[T]
}

func (t *RBTree[T]) Rotate(oldSubRoot *RBNode[T], direction NodeSide) *RBNode[T] {
	parent := oldSubRoot.Parent
	parentSide := oldSubRoot.Side()
	newSubRoot := oldSubRoot.Children[1-direction]
	if newSubRoot == nil {
		return oldSubRoot
	}
	swapChild := newSubRoot.Children[direction]
	oldSubRoot.Children[1-direction] = swapChild
	if swapChild != nil {
		swapChild.Parent = oldSubRoot
	}
	newSubRoot.Children[direction] = oldSubRoot
	oldSubRoot.Parent = newSubRoot
	newSubRoot.Parent = parent
	if parent != nil {
		parent.Children[parentSide] = newSubRoot
	} else {
		t.Root = newSubRoot
	}
	return newSubRoot
}

func (t *RBTree[T]) InOrderSuccessor(node *RBNode[T]) *RBNode[T] {
	succ := node.Children[Right]
	if succ == nil {
		return nil
	}
	for succ.Children[Left] != nil {
		succ = succ.Children[Left]
	}
	return succ
}

func (t *RBTree[T]) InOrderPredecessor(node *RBNode[T]) *RBNode[T] {
	pred := node.Children[Left]
	if pred == nil {
		return nil
	}
	for pred.Children[Right] != nil {
		pred = pred.Children[Right]
	}
	return pred
}

// Adapted from: https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
func (t *RBTree[T]) Insert(node *RBNode[T], parent *RBNode[T], side NodeSide) {
	node.Color = Red
	node.Children = [2]*RBNode[T]{nil, nil}
	node.Parent = parent
	var grandParent, uncle *RBNode[T]
	if parent == nil {
		t.Root = node
		return
	}
	parent.Children[side] = node
	keepGoing := true
	for keepGoing {
		if parent.IsBlack() {
			return
		}
		grandParent = parent.Parent
		if grandParent == nil {
			goto CASE_I4
		}
		side = parent.Side()
		uncle = node.Uncle()
		if uncle == nil || uncle.IsBlack() {
			goto CASE_56
		}
		parent.Color = Black
		uncle.Color = Black
		grandParent.Color = Red
		node = grandParent
		parent = node.Parent
		keepGoing = parent != nil
	}
	return
CASE_I4:
	parent.Color = Black
	return
CASE_56:
	if node == parent.Children[1-side] {
		t.Rotate(parent, side)
		node = parent
		parent = grandParent.Children[side]
	}
	t.Rotate(grandParent, 1-side)
	parent.Color = Black
	grandParent.Color = Red
	return
}

// Adapted from: https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
func (t *RBTree[T]) Delete(node *RBNode[T]) {
	if node.Parent == nil && node.Children[Left] == nil && node.Children[Right] == nil {
		t.Root = nil
		node = nil
		return
	}
	if node.Children[Left] != nil && node.Children[Right] != nil {
		successor := t.InOrderSuccessor(node)
		node.SwapWith(successor)
	}
	if node.IsRed() {
		node.Parent.Children[node.Side()] = nil
		node = nil
		return
	}
	if leftChild, rightChild := node.Children[Left] != nil && node.Children[Right] == nil, node.Children[Left] == nil && node.Children[Right] != nil; leftChild || rightChild {
		var side NodeSide
		if rightChild {
			side = Right
		}
		node.Children[side].Color = Black
		node.SwapWith(node.Children[side])
		node.Parent.Children[node.Side()] = nil
		node = nil
		return
	}
	parent := node.Parent
	side := node.Side()
	parent.Children[side] = nil
	var sibling, closeNephew, distNephew *RBNode[T]
	keepGoing := true
	for keepGoing {
		sibling = parent.Children[1-side]
		distNephew = sibling.Children[1-side]
		closeNephew = sibling.Children[side]
		if sibling.IsRed() {
			goto CASE_D3
		}
		if distNephew != nil && distNephew.IsRed() {
			goto CASE_D6
		}
		if closeNephew != nil && closeNephew.IsRed() {
			goto CASE_D5
		}
		if parent.IsRed() {
			goto CASE_D4
		}
		sibling.Color = Red
		node = parent
		parent = node.Parent
		keepGoing = parent != nil
		if keepGoing {
			side = node.Side()
		}
	}
	return
CASE_D3:
	t.Rotate(parent, side)
	parent.Color = Red
	sibling.Color = Black
	sibling = closeNephew
	distNephew = sibling.Children[1-side]
	if distNephew != nil && distNephew.IsRed() {
		goto CASE_D6
	}
	if closeNephew != nil && closeNephew.IsRed() {
		goto CASE_D5
	}
CASE_D4:
	sibling.Color = Red
	parent.Color = Black
	return
CASE_D5:
	t.Rotate(sibling, 1-side)
	sibling.Color = Red
	closeNephew.Color = Black
	distNephew = sibling
	sibling = closeNephew
CASE_D6:
	t.Rotate(parent, side)
	sibling.Color = parent.Color
	parent.Color = Black
	distNephew.Color = Black
	return
}
