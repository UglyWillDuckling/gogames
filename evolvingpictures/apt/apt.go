package apt

import (
	"math"
	"math/rand"
	"strconv"
	"vlado/game/noise"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
	AddRandom(node Node)
}

type LeafNode struct{}

func (leaf *LeafNode) AddRandom(node Node) {
	panic("ERROR: Your tried to add a node to a leaf node")
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

type OpX LeafNode

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY LeafNode

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

func (op *OpConstant) OpConstant() string {
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
