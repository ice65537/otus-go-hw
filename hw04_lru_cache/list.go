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
}

func NewList() List {
	return new(list)
}

type list struct {
	front *ListItem
	back  *ListItem
	count int
}

func (l list) Len() int {
	return l.count
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	x := ListItem{Value: v, Next: l.front, Prev: nil}
	if l.front != nil {
		(*l.front).Prev = &x
	}
	l.front = &x
	if l.back == nil {
		l.back = &x
	}
	l.count++
	return &x
}

func (l *list) PushBack(v interface{}) *ListItem {
	x := ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		(*l.back).Next = &x
	}
	l.back = &x
	if l.front == nil {
		l.front = &x
	}
	l.count++
	return &x
}

func (l *list) Remove(i *ListItem) {
	if l.count > 0 {
		if i == l.front {
			l.front = i.Next
		}
		if i == l.back {
			l.back = i.Prev
		}
		if i.Next != nil {
			(*i.Next).Prev = i.Prev
		}
		if i.Prev != nil {
			(*i.Prev).Next = i.Next
		}
		l.count--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}
