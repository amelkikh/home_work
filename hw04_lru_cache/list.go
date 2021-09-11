package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
	list  *list
}

type list struct {
	head, tail *ListItem
	count      int
}

func (l list) Len() int {
	return l.count
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.head,
		list:  l,
	}
	if l.head != nil {
		l.head.Prev = item
	} else {
		l.tail = item
	}
	l.head = item
	l.count++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.tail,
		list:  l,
	}
	if l.head == nil {
		l.head = item
	}
	if l.tail != nil {
		l.tail.Next = item
	}
	l.tail = item
	l.count++

	return item
}

func (l *list) Remove(item *ListItem) {
	if item == nil || item.list == nil || item.list != l {
		return
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.head = item.Next
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.tail = item.Prev
	}
	// avoid memory leaks
	item.Next = nil
	item.Prev = nil
	// unlink item to avoid incorrect behaviour with other operations with list
	item.list = nil
	l.count--
}

func (l *list) MoveToFront(item *ListItem) {
	if item.Prev == nil || item.list == nil {
		// Don't need to move front first item or if item unlinked (deleted) from list
		return
	}
	l.Remove(item)
	l.head.Prev = item
	item.Next = l.head
	l.head = item
	item.list = l
	l.count++
}

func NewList() List {
	return new(list)
}
