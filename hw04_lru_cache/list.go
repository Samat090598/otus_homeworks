package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Init()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v}

	if l.head == nil {
		l.head = &item
		l.tail = &item
	} else {
		item.Next = l.head
		l.head.Prev = &item
		l.head = &item
	}

	l.length++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v}

	if l.tail == nil {
		l.head = &item
		l.tail = &item
	} else {
		item.Prev = l.tail
		l.tail.Next = &item
		l.tail = &item
	}

	l.length++
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		l.tail = i.Prev
	}

	i.Prev.Next = i.Next
	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}

func (l *list) Init() {
	l.head = nil
	l.tail = nil
	l.length = 0
}
