package async

import (
  "container/list"
  "sync"
)

/*
  Used to contain the Routine functions to be processed
*/
type List struct {
  *list.List

  Wait sync.WaitGroup
}

/*
  Create a new list
*/
func New() *List {
  return &List{
    List: list.New(),
  }
}

/*
  Add a Routine function to the current list
*/
func (l *List) Add(routine Routine) (*List, *list.Element) {
  element := l.PushBack(routine)
  return l, element
}

/*
  Add multiple Routine functions to the current list
*/
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
  Remove an element from the current list
*/
func (l *List) Remove(element *list.Element) (*List, Routine) {
  routine := l.List.Remove(element).(Routine)
  return l, routine
}
