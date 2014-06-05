package async

import "testing"

func TestMapString(t *testing.T) {
  str := []string{
    "test",
    "test2",
    "test3",
    "test4",
    "test5",
  }

  expects := []string{
    "test1",
    "test2",
    "test3",
    "test4",
    "test5",
  }

  mapper := func(done Done, args ...interface{}) {
    Status("Hit string")
    Status("Args: %+v\n", args)
    if args[1] == 0 {
      done(nil, "test1")
      return
    }
    done(nil, args[0])
  }

  final := func(err error, results ...interface{}) {
    Status("Hit string end")
    Status("Results: %+v\n", results)
    for i := 0; i < len(results); i++ {
      if results[i] != expects[i] {
        t.Errorf("Did not map correctly.")
        break
      }
    }
  }

  Map(str, mapper, final)
}

func TestMapInt(t *testing.T) {
  ints := []int{1, 2, 3, 4, 5}

  expects := []int{2, 4, 6, 8, 10}

  mapper := func(done Done, args ...interface{}) {
    Status("Hit int")
    Status("Args: %+v\n", args)
    done(nil, args[0].(int)*2)
  }

  final := func(err error, results ...interface{}) {
    Status("Hit int end")
    Status("Results: %+v\n", results)
    for i := 0; i < len(results); i++ {
      if results[i] != expects[i] {
        t.Errorf("Did not map correctly.")
        break
      }
    }
  }

  Map(ints, mapper, final)
}

func TestMapBool(t *testing.T) {
  bools := []bool{true, false, false, true, false}

  expects := []bool{true, true, false, true, false}

  mapper := func(done Done, args ...interface{}) {
    Status("Hit bool")
    Status("Args: %+v\n", args)
    if args[1] == 1 {
      done(nil, true)
      return
    }
    done(nil, args[0])
  }

  final := func(err error, results ...interface{}) {
    Status("Hit bool end")
    Status("Results: %+v\n", results)
    for i := 0; i < len(results); i++ {
      if results[i] != expects[i] {
        t.Errorf("Did not map correctly.")
        break
      }
    }
  }

  Map(bools, mapper, final)
}
