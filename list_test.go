package async

import "testing"

func TestAdd(t *testing.T) {
  Status("Creating list")
  list := New()

  Status("Adding routine")
  list.Add(func(done Done, args ...interface{}) {
    Status("Args: %+v", args)
  })

  if list.Len() > 0 {
    Status("Added routine to list")
  }
}

func TestMultiple(t *testing.T) {
  Status("Creating list")
  list := New()

  Status("Adding multiple routines")
  list, _ = list.Multiple(func(done Done, args ...interface{}) {
    Status("Args: %+v", args)
  }, func(done Done, args ...interface{}) {
    Status("Args2: %+v", args)
  })

  if list.Len() == 2 {
    Status("Added multiple routines to list")
  }
}

func TestRemove(t *testing.T) {
  Status("Creating list")
  list := New()

  Status("Adding routine")
  list, elem := list.Add(func(done Done, args ...interface{}) {
    Status("Args: %+v", args)
  })

  if list.Len() > 0 {
    Status("Added routine to list")
  }

  Status("Removing element")
  list.Remove(elem)

  if list.Len() == 0 {
    Status("Removed element from list")
  }
}
