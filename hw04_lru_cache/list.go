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

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}

	return l.front
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}

	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	var nli *ListItem

	if front := l.Front(); front != nil {
		nli = &ListItem{
			Prev:  nil,
			Next:  front,
			Value: v,
		}

		front.Prev = nli
		l.front = nli

		l.len++
	} else {
		nli = &ListItem{
			Prev:  nil,
			Next:  nil,
			Value: v,
		}

		l.front = nli
		l.back = nli

		l.len++
	}

	return nli
}

func (l *list) PushBack(v interface{}) *ListItem {
	var nli *ListItem

	if back := l.Back(); back != nil {
		nli = &ListItem{
			Prev:  back,
			Next:  nil,
			Value: v,
		}

		back.Next = nli
		l.back = nli

		l.len++
	} else {
		nli = &ListItem{
			Prev:  nil,
			Next:  nil,
			Value: v,
		}

		l.front = nli
		l.back = nli

		l.len++
	}

	return nli
}

func (l *list) Remove(i *ListItem) {
	if i != nil {
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}

		if i.Next != nil {
			i.Next.Prev = i.Prev
		}

		i.Prev = nil
		i.Next = nil

		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i != nil && i.Prev != nil {
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}

		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.back = i.Prev
		}

		i.Prev = nil
		i.Next = l.front

		l.front.Prev = i
		l.front = i
	}
}

func NewList() List {
	return new(list)
}
