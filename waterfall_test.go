package async_test

import (
  "fmt"
  "github.com/Southern/async"
  "testing"
  "time"
)

func TestWaterfall(t *testing.T) {
  Status("Calling Waterfall")
  async.Waterfall([]async.Routine{
    func(done async.Done, args ...interface{}) {
      Status("First waterfall function")
      Status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done async.Done, args ...interface{}) {
      Status("Second waterfall function")
      Status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg4", "arg5", "arg6")
    },

    func(done async.Done, args ...interface{}) {
      Status("Third waterfall function")
      Status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg7", "arg8", "arg9")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      t.Errorf("Waterfall threw an unexpected error: %+v", err)
      return
    }

    Status("Waterfall completed with results: %+v", results)
  })
}

func TestWaterfallError(t *testing.T) {
  Status("Calling Waterfall")
  async.Waterfall([]async.Routine{
    func(done async.Done, args ...interface{}) {
      Status("First waterfall function")
      Status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(nil, "arg1", "arg2", "arg3")
    },

    func(done async.Done, args ...interface{}) {
      Status("Second waterfall function")
      Status("Called with arguments: %+v", args)
      time.Sleep(time.Second)
      done(fmt.Errorf("Test error"))
    },

    func(done async.Done, args ...interface{}) {
      t.Errorf("The second waterfall function did not stop when it errored.")
    },
  }, func(err error, results ...interface{}) {
    if err != nil {
      Status("Waterfall exited with error: %+v", err)
      return
    }

    t.Errorf("Waterfall did not throw an error as expected")
  })
}
