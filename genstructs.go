package genstructs

type NodeSide uint8

const (
	Left  NodeSide = 0
	Right NodeSide = 1
)

type SplitFunc[T any] func(value T) (left T, right T)

type MatchFunc[T any] func(value T, existing T) (match bool, sideNoMatch NodeSide)

type CompareFunc[T any] func(value, existing T) NodeSide
