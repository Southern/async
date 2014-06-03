package mapper

import (
  "fmt"
  "git.aviuslabs.net/golang/async"
  "testing"
)

func TestMapString(t *testing.T) {
  str := []string{
    "test",
    "test2",
    "test3",
    "test4",
    "test5",
  }

  mapper := func(done async.Done, args ...interface{}) {
    println("Hit string")
    fmt.Printf("Args: %+v\n", args)
    done(nil, args...)
  }

  final := func(err error, results ...interface{}) {
    println("Hit string end")
    fmt.Printf("Results: %+v\n", results)
  }

  Map(str, mapper, final)
}

func TestMapInt(t *testing.T) {
  ints := []int{
    1, 2, 3, 4, 5,
  }

  mapper := func(done async.Done, args ...interface{}) {
    println("Hit bool")
    fmt.Printf("Args: %+v\n", args)
    done(nil, args...)
  }

  final := func(err error, results ...interface{}) {
    println("Hit string end")
    fmt.Printf("Results: %+v\n", results)
  }

  Map(ints, mapper, final)
}

func TestMapBool(t *testing.T) {
  bools := []bool{
    true, false, false, true, false,
  }

  mapper := func(done async.Done, args ...interface{}) {
    println("Hit bool")
    fmt.Printf("Args: %+v\n", args)
    done(nil, args...)
  }

  final := func(err error, results ...interface{}) {
    println("Hit string end")
    fmt.Printf("Results: %+v\n", results)
  }

  Map(bools, mapper, final)
}
