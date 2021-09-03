package linkedlist

import (
	"fmt"
	"testing"
)

func TestList_Push(t *testing.T) {
	l := NewList()

	l.Push(&Item{Value: 1})
	l.Push(&Item{Value: 2})
	fmt.Println(*l.head.Next, l.Len())
	l.Remove(l.head.Next)
	fmt.Println(*l.head.Next, l.Len())
	l.Remove(l.head.Next)
	fmt.Println(*l.head, l.Len())
}
