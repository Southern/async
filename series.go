package async

/*
  Shorthand to List.RunSeries without having to manually create a new
  list, add the routines, etc.
*/
func Series(routines []Routine, callbacks ...Done) {
  l := New()
  l.Multiple(routines...)

  l.RunSeries(callbacks...)
}

/*
  Run all of the Routine functions in a series effect.

  If there is an error, series will immediately exit and trigger the
  callbacks with the error.

  There are no arguments passed between the routines that are used in series.
  It is just for commands that need to run asynchronously without seeing the
  results of its previous routine.

  For example, take a look at one of the tests for this function:
    func TestSeries(t *testing.T) {
      counter := 0

      Status("Calling Series")
      async.Series([]async.Routine{
        func(done async.Done, args ...interface{}) {
          Status("Increasing counter...")
          counter++
          done(nil)
        },
        func(done async.Done, args ...interface{}) {
          Status("Increasing counter...")
          counter++
          done(nil)
        },
        func(done async.Done, args ...interface{}) {
          Status("Increasing counter...")
          counter++
          done(nil)
        },
        func(done async.Done, args ...interface{}) {
          Status("Increasing counter...")
          counter++
          done(nil)
        },
      }, func(err error, results ...interface{}) {
        if err != nil {
          t.Errorf("Unexpected error: %s", err)
          return
        }

        if counter != 4 {
          t.Errorf("Not all routines were completed.")
          return
        }

        Status("Counter: %d", counter)
      })
    }

*/
func (l *List) RunSeries(callbacks ...Done) {
  fall := fallSeries(l, callbacks...)
  next := nextSeries(l, callbacks...)

  l.Wait.Add(l.Len())

  fall(next)
}

func fallSeries(l *List, callbacks ...Done) func(Done, ...interface{}) {
  return func(next Done, args ...interface{}) {
    e := l.Front()
    _, r := l.Remove(e)

    // Run the first series routine and give it the next function, and
    // any arguments that were provided
    go r(next)
    l.Wait.Wait()
  }
}

func nextSeries(l *List, callbacks ...Done) Done {
  fall := fallSeries(l, callbacks...)

  return func(err error, args ...interface{}) {
    next := nextSeries(l, callbacks...)

    l.Wait.Done()
    if err != nil || l.Len() == 0 {
      // Just in case it's an error, let's make sure we've cleared
      // all of the sync.WaitGroup waits that we initiated.
      for i := 0; i < l.Len(); i++ {
        l.Wait.Done()
      }

      // Send the results to the callbacks
      for i := 0; i < len(callbacks); i++ {
        callbacks[i](err)
      }
      return
    }

    // Run the next series routine with any arguments that were provided
    fall(next)
    return
  }
}