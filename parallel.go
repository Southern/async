package async

func Parallel(routines []Routine, callbacks ...Done) {
  l := New()
  l.Multiple(routines...)

  l.RunParallel(callbacks...)
}

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

  l.Wait.Add(l.Len())

  go func() {
    for {
      r := <-result

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

    go func() {
      r(func(err error, args ...interface{}) {
        if _error != nil {
          return
        }

        if err != nil {
          result <- err
          return
        }

        result <- args
      })

      l.Wait.Done()
    }()
  }

  l.Wait.Wait()

  close(result)

  if _error == nil {
    final(nil, results...)
  }
}
