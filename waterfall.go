package async

func Waterfall(routines []Routine, callbacks ...Done) {
  l := New()
  l.Multiple(routines...)

  l.RunWaterfall(callbacks...)
}

func (l *List) RunWaterfall(callbacks ...Done) {
  fall := fall(l, callbacks...)
  next := next(l, callbacks...)

  l.Wait.Add(l.Len())

  fall(next)
}

func fall(l *List, callbacks ...Done) func(Done, ...interface{}) {
  return func(next Done, args ...interface{}) {
    e := l.Front()
    _, r := l.Remove(e)

    // Run the first waterfall routine and give it the next function, and
    // any arguments that were provided
    r(next, args...)
    l.Wait.Wait()
  }
}

func next(l *List, callbacks ...Done) Done {
  fall := fall(l, callbacks...)

  return func(err error, args ...interface{}) {
    next := next(l, callbacks...)

    l.Wait.Done()
    if err != nil || l.Len() == 0 {
      // Just in case it's an error, let's make sure we've cleared
      // all of the sync.WaitGroup waits that we initiated.
      for i := 0; i < l.Len(); i++ {
        l.Wait.Done()
      }

      // Send the results to the callbacks
      for i := 0; i < len(callbacks); i++ {
        callbacks[i](err, args...)
      }
      return
    }

    // Run the next waterfall routine with any arguments that were provided
    fall(next, args...)
    return
  }
}
