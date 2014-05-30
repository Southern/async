package async

import (
  "fmt"
  "strings"
  "testing"
  "time"
)

func status(format string, args ...interface{}) {
  if testing.Verbose() {
    fmt.Printf(strings.TrimSpace(format)+"\n", args...)
  }
}

func TestWaterfall(t *testing.T) {
  status("Calling Waterfall")
  Waterfall([]Routine{
    func(done Done, args ...interface{}) {
      status("First waterfall function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      status("Second waterfall function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg4", "arg5", "arg6")
    },

    func(done Done, args ...interface{}) {
      status("Third waterfall function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg7", "arg8", "arg9")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      t.Errorf("Waterfall threw an unexpected error: %+v", err)
      return
    }

    status("Waterfall completed with results: %+v", results)
  })
}

func TestWaterfallError(t *testing.T) {
  status("Calling Waterfall")
  Waterfall([]Routine{
    func(done Done, args ...interface{}) {
      status("First waterfall function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      status("Second waterfall function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(fmt.Errorf("Test error"))
    },

    func(done Done, args ...interface{}) {
      t.Errorf("The second waterfall function did not stop when it errored.")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      status("Waterfall exited with error: %+v", err)
      return
    }

    t.Errorf("Waterfall did not throw an error as expected")
  })
}

func TestClearWaterfall(t *testing.T) {
  status("Creating new list")
  list := New()

  status("Adding test routine")
  err, list := list.Add(func(done Done, args ...interface{}) {
    t.Errorf("Test")
  })

  status("Checking to see if adding returned an error")
  if err != nil {
    t.Errorf("Got an error: %s", err)
  }

  status("Clearing the list")
  list.Clear()

  status("Attempting to run an empty list")
  list.Run(func(err error, results ...interface{}) {
    if err == nil {
      t.Error("Expected an error")
      return
    }

    if err == NoRoutines {
      status("List was cleared")
      return
    }

    t.Errorf("Unknown error: %+v", err)
  })
}

func TestParallel(t *testing.T) {
  status("Calling parallel")
  Parallel([]Routine{
    func(done Done, args ...interface{}) {
      status("First parallel function")
      status("Called with arguments: %+v", args)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      status("Second parallel function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg4", "arg5", "arg6")
    },

    func(done Done, args ...interface{}) {
      status("Third parallel function")
      status("Called with arguments: %+v", args)
      done(nil, "arg7", "arg8", "arg9")
    },
  }, func(err error, results ...interface{}) {
    status("Parallel err: %+v", err)
    status("Parallel results: %+v", results)

    if err != nil {
      t.Errorf("Parallel threw an unexpected error: %+v", err)
      return
    }

    status("Parallel completed with results: %+v", results)
  })
}

func TestParallelError(t *testing.T) {
  status("Calling Parallel")
  Parallel([]Routine{
    func(done Done, args ...interface{}) {
      status("First parallel function")
      status("Called with arguments: %+v", args)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      status("Second parallel function")
      status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(fmt.Errorf("Test error"))
    },

    func(done Done, args ...interface{}) {
      status("Third parallel function")
      status("Called with arguments: %+v", args)
      done(nil, "arg4", "arg5", "arg6")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      status("Parallel exited with error: %+v", err)
      return
    }

    t.Errorf("Parallel did not throw an error as expected")
  })
}

func TestClearParallel(t *testing.T) {
  status("Creating new list")
  list := New()

  status("Adding test routine")
  err, list := list.Add(func(done Done, args ...interface{}) {
    t.Errorf("Test")
  })

  status("Checking to see if adding returned an error")
  if err != nil {
    t.Errorf("Got an error: %s", err)
  }

  status("Clearing the list")
  list.Clear()

  status("Attempting to run an empty list")
  list.RunParallel(func(err error, results ...interface{}) {
    if err == nil {
      t.Error("Expected an error")
      return
    }

    if err == NoRoutines {
      status("List was cleared")
      return
    }

    t.Errorf("Unknown error: %+v", err)
  })
}

func ExampleParallel() {
  Parallel([]Routine{
    func(done Done, args ...interface{}) {
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      done(nil, "arg4", "arg5", "arg6")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      fmt.Errorf("Error: %+v\n", err)
      return
    }

    fmt.Printf("Results: %+v\n", results)
  })

  // Output:
  // Results: [arg4 arg5 arg6 arg1 arg2 arg3]
}

func ExampleWaterfall() {
  Waterfall([]Routine{
    func(done Done, args ...interface{}) {
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done Done, args ...interface{}) {
      time.Sleep(time.Second)
      done(nil, "arg4", "arg5", "arg6")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      fmt.Errorf("Error: %+v\n", err)
      return
    }

    fmt.Printf("Results: %+v\n", results)
  })

  // Output:
  // Results: [arg4 arg5 arg6]
}
