package async

/*

Parallel is a shorthand function to List.RunParallel without having to
manually create a new list, add the routines, etc.

*/
func Parallel(routines []Routine, callbacks ...Done) {
  l := New()
  l.Multiple(routines...)

  l.RunParallel(callbacks...)
}

/*

RunParallel will run all of the Routine functions from the current list in
parallel mode.

All of the arguments returned in a Routine's Done function will be combined
and returned in the callbacks that are provided.

If there is an error, any further results will be discarded but it will not
immediately exit. It will continue to run all of the other Routine functions
that were passed into it. This is because by the time the error is sent, the
goroutines have already been started. At this current time, there is no way
to cancel a sleep timer in Go.

For example:
  async.Parallel([]async.Routine{
    func(done async.Done, args ...interface{}) {
      time.Sleep(20 * time.Second)
      done(nil, "Won't trigger the callbacks because error has been sent")
    },
    func(done async.Done, args ...interface{}) {
      done(fmt.Errorf("Test error"))
    }
  }, func(err error, results ...interface{}) {
    if err != nil {
      fmt.Printf("Error: %s", err)
      return
    }

    fmt.Printf("Args: %s", args)
  })

If you were to run this example, you would see the error happen immediately.
However, you would also notice that the program doesn't immediately exit.
That is because it is still waiting for responses that it silently discards,
since an error has already occurred.

*/
func (l *List) RunParallel(callbacks ...Done) {
  var (
    results = make([]interface{}, 0)

    result = make(chan interface{})

    _error error

    final = func(err error, results ...interface{}) {
      for i := 0; i < len(callbacks); i++ {
        if err != nil {
          callbacks[i](err)
        } else {
          callbacks[i](err, results...)
        }
      }
    }
  )

  defer close(result)

  l.Wait.Add(l.Len())

  go func() {
    for {
      r := <-result

      if _error != nil {
        continue
      }

      switch r.(type) {
      case error:
        _error = r.(error)
        final(_error)

      case []interface{}:
        results = append(results, r.([]interface{})...)
      }
    }
  }()

  for l.Len() > 0 {
    e := l.Front()
    _, r := l.Remove(e)

    go r(func(err error, args ...interface{}) {
      defer l.Wait.Done()

      if err != nil {
        result <- err
        return
      }

      result <- args
    })
  }

  l.Wait.Wait()

  if _error == nil {
    final(nil, results...)
  }
}
