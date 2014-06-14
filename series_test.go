package async_test

import (
  "fmt"
  "github.com/Southern/async"
  "testing"
)

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

func TestSeriesError(t *testing.T) {
  counter := 0

  Status("Calling Series")
  async.Series([]async.Routine{
    func(done async.Done, args ...interface{}) {
      Status("Increasing counter...")
      counter++
      done(nil)
    },
    func(done async.Done, args ...interface{}) {
      Status("Sending error...")
      done(fmt.Errorf("Test error"))
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
    if err == nil {
      t.Errorf("Did not throw an error as expected")
      return
    }

    Status("Got error: %s", err)
  })
}

func TestSeriesParallel(t *testing.T) {
  counter := 0

  Status("Calling Series")
  async.SeriesParallel([]async.Routine{
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

func TestSeriesParallelError(t *testing.T) {
  counter := 0

  Status("Calling Series")
  async.SeriesParallel([]async.Routine{
    func(done async.Done, args ...interface{}) {
      Status("Increasing counter...")
      counter++
      done(nil)
    },
    func(done async.Done, args ...interface{}) {
      Status("Sending error...")
      done(fmt.Errorf("Test error"))
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
    if err == nil {
      t.Errorf("Did not throw an error as expected")
      return
    }

    Status("Got error: %s", err)
  })
}
