// Package mytype
//  Forked to
//  https://github.com/floyernick/Data-Structures-and-Algorithms
package mytype

type node struct {
	data Entity
	next *node
}

type EntityQueue struct {
	rear *node
}

func (list *EntityQueue) Enqueue(i Entity) {
	data := &node{data: i}
	if list.rear != nil {
		data.next = list.rear
	}
	list.rear = data
}

func (list *EntityQueue) Dequeue() (Entity, bool) {
	if list.rear == nil {
		return Entity{}, false
	}
	if list.rear.next == nil {
		i := list.rear.data
		list.rear = nil
		return i, true
	}
	current := list.rear
	for {
		if current.next.next == nil {
			i := current.next.data
			current.next = nil
			return i, true
		}
		current = current.next
	}
}

func (list *EntityQueue) Peek() (Entity, bool) {
	if list.rear == nil {
		return Entity{}, false
	}
	return list.rear.data, true
}

func (list *EntityQueue) Get() []Entity {
	var items []Entity
	current := list.rear
	for current != nil {
		items = append(items, current.data)
		current = current.next
	}
	return items
}

func (list *EntityQueue) IsEmpty() bool {
	return list.rear == nil
}

func (list *EntityQueue) Clear() {
	list.rear = nil
}
