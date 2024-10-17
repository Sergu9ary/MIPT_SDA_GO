//go:build !solution

package treeiter

type Tree[T any] interface {
	Left() *T
	Right() *T
}

func DoInOrder[T Tree[T]](root *T, process func(*T)) {
	if root == nil {
		return
	}
	DoInOrder((*root).Left(), process)
	process(root)
	DoInOrder((*root).Right(), process)
}
