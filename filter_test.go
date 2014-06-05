package async

import (
  "fmt"
  "testing"
)

func TestFilter(t *testing.T) {
  str := []string{
    "test1",
    "test2",
    "test3",
    "test4",
    "test5",
  }

  mapper := func(done Done, args ...interface{}) {
    println("Hit string")
    fmt.Printf("Args: %+v\n", args)
    if args[0] == "test3" {
      done(nil, false)
      return
    }
    done(nil, true)
  }

  final := func(err error, results ...interface{}) {
    println("Hit string end")
    fmt.Printf("Results: %+v\n", results)
    if results[2] != "test4" {
      t.Errorf("Did not filter correctly.")
    }
  }

  Filter(str, mapper, final)
}
