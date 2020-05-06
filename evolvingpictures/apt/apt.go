package apt

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"vlado/game/noise"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
	AddRandom(node Node)
	NodeCounts() (nodeCount int, nilCount int)
}

type LeafNode struct{}

func (leaf *LeafNode) AddRandom(node Node) {
	// panic("ERROR: Your tried to add a node to a leaf node")
	fmt.Println("ERROR: Your tried to add a node to a leaf node")
}

func (leaf *LeafNode) NodeCounts() (nodeCount, nilCount int) {
	return 1, 0
}

type SingleNode struct {
	Child Node
}

func (singleNode *SingleNode) AddRandom(node Node) {
	if singleNode.Child == nil {
		singleNode.Child = node
	} else {
		singleNode.Child.AddRandom(node)
	}
}

func (single *SingleNode) NodeCounts() (nodeCount, nilCount int) {
	if single.Child == nil {
		return 1, 1
	}

	childCount, childNilCount := single.Child.NodeCounts()
	return 1 + childCount, childNilCount
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

func (doubleNode *DoubleNode) AddRandom(node Node) {
	r := rand.Intn(2)

	if r == 0 {
		if doubleNode.LeftChild == nil {
			doubleNode.LeftChild = node
		} else {
			doubleNode.LeftChild.AddRandom(node)
		}
	} else {
		if doubleNode.RightChild == nil {
			doubleNode.RightChild = node
		} else {
			doubleNode.RightChild.AddRandom(node)
		}
	}
}

func (doubleNode *DoubleNode) NodeCounts() (nodeCount, nilCount int) {
	var leftCount, leftNilCount, rightCount, rightNilCount int

	if doubleNode.LeftChild == nil {
		leftNilCount = 1
		leftCount = 0
	} else {
		leftCount, leftNilCount = doubleNode.LeftChild.NodeCounts()
	}

	if doubleNode.RightChild == nil {
		rightNilCount = 1
		rightCount = 0
	} else {
		rightCount, rightNilCount = doubleNode.RightChild.NodeCounts()
	}

	return 1 + leftCount + rightCount, leftNilCount + rightNilCount
}

type OpSin struct {
	SingleNode
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

func (op *OpSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OpCos struct {
	SingleNode
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Child.Eval(x, y))))
}

func (op *OpCos) String() string {
	return "( Cos " + op.Child.String() + " )"
}

type OpPlus struct {
	DoubleNode
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

func (op *OpPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMinus struct {
	DoubleNode
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) - op.RightChild.Eval(x, y)
}

func (op *OpMinus) String() string {
	return "( - " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMult struct {
	DoubleNode
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) * op.RightChild.Eval(x, y)
}

func (op *OpMult) String() string {
	return "( * " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpDiv struct {
	DoubleNode
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) / op.RightChild.Eval(x, y)
}

func (op *OpDiv) String() string {
	return "( / " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type opAtan struct {
	SingleNode
}

func (op *opAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Child.Eval(x, y))))
}

func (op *opAtan) String() string {
	return "( Atan " + op.Child.String() + " )"
}

type opAtan2 struct {
	DoubleNode
}

func (op *opAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(op.LeftChild.Eval(x, y)), float64(op.RightChild.Eval(x, y))))
}

func (op *opAtan2) String() string {
	return "( Atan2 " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpNoise struct {
	DoubleNode
}

func (op *OpNoise) Eval(x, y float32) float32 {
	return 80*noise.Snoise2(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y)) - 2.0
}

func (op *OpNoise) String() string {
	return "( Noise " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpLog2 struct {
	SingleNode
}

func (op *OpLog2) Eval(x, y float32) float32 {
	return float32(math.Log2(float64(op.Child.Eval(x, y))))
}

func (op *OpLog2) String() string {
	return "( Log2 " + op.Child.String() + " )"
}

type OpX struct {
	LeafNode
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY struct {
	LeafNode
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}

type OpConstant struct {
	LeafNode
	value float32
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 9, 32)
}

func GetRandomNode() Node {
	r := rand.Intn(9)

	switch r {
	case 0:
		return &OpPlus{}
	case 1:
		return &OpMinus{}
	case 2:
		return &OpMult{}
	case 3:
		return &OpDiv{}
	case 4:
		return &opAtan2{}
	case 5:
		return &opAtan{}
	case 6:
		return &OpCos{}
	case 7:
		return &OpSin{}
	case 8:
		return &OpNoise{}
	}

	return &OpNoise{}
}

func GetRandomLeaf() Node {
	r := rand.Intn(3)

	switch r {
	case 0:
		return &OpX{}
	case 1:
		return &OpY{}
	case 2:
		return &OpConstant{LeafNode{}, rand.Float32()*2 - 1}
	}

	panic("Error in get random Leaf")
}
