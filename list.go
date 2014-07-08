package async

import (
	"container/list"
	"sync"
)

/*

List is used to contain the Routine functions to be processed

This list inherits https://godoc.org/container/list and contains all
of the functionality that it contains, with a minor tweak to Remove. Instead
of Remove returning the element, it returns our routine. This is used to
ensure that our Routine is removed from the list before it's ran, and
therefore isn't able to be called again.

*/
type List struct {
	*list.List

	Wait sync.WaitGroup
}

// New will create a new List instance
func New() *List {
	return &List{
		List: list.New(),
	}
}

// Add will add a single Routine function to the current list
func (l *List) Add(routine Routine) (*List, *list.Element) {
	element := l.PushBack(routine)
	return l, element
}

// Multiple will add multiple Routine functions to the current list
func (l *List) Multiple(routines ...Routine) (*List, []*list.Element) {
	var (
		elements = make([]*list.Element, 0)
	)

	for i := 0; i < len(routines); i++ {
		_, e := l.Add(routines[i])
		elements = append(elements, e)
	}

	return l, elements
}

/*

Remove deletes a Routine element from the current list

An example of removing a Routine would be:
  err, routine := l.Remove(l.Front())

*/
func (l *List) Remove(element *list.Element) (*List, Routine) {
	routine := l.List.Remove(element).(Routine)
	return l, routine
}
