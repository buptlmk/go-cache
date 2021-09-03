package linkedlist

import (
	"sync"
)

// 线程安全的双向链表

type Item struct {
	Value interface{}

	Prev *Item
	Next *Item
}

type List struct {
	head *Item
	tail *Item

	number int

	lock sync.RWMutex
}

func NewList() *List {

	head := &Item{}

	return &List{
		head: head,
		tail: head,
	}
}

func (l *List) Push(value *Item) {
	l.lock.Lock()
	defer l.lock.Unlock()
	temp := l.tail
	temp.Next = value
	value.Prev = temp
	l.tail = value
	l.number++
}

func (l *List) Remove(value *Item) {
	if value == nil || value.Prev == nil {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	temp := value.Prev
	temp.Next = value.Next

	if value.Next != nil {
		value.Next.Prev = temp
	} else {
		l.tail = temp
	}
	l.number--
}

func (l *List) Front() (value *Item) {

	return l.head.Next
}

func (l *List) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.number
}
