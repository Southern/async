package async

import (
	"reflect"
)

/*

Filter allows you to filter out information from a slice in Waterfall mode.

You must call the Done function with false as its first argument if you do
not want the data to be present in the results. No other arguments will
affect the performance of this function. When calling the Done function,
an error will cause the filtering to immediately exit.

For example, take a look at one of the tests for this function:
  func TestFilterString(t *testing.T) {
    str := []string{
      "test1",
      "test2",
      "test3",
      "test4",
      "test5",
    }

    expects := []string{
      "test1",
      "test2",
      "test4",
      "test5",
    }

    mapper := func(done async.Done, args ...interface{}) {
      Status("Hit string")
      Status("Args: %+v\n", args)
      if args[0] == "test3" {
        done(nil, false)
        return
      }
      done(nil, true)
    }

    final := func(err error, results ...interface{}) {
      Status("Hit string end")
      Status("Results: %+v\n", results)
      for i := 0; i < len(results); i++ {
        if results[i] != expects[i] {
          t.Errorf("Did not filter correctly.")
          break
        }
      }
    }

    async.Filter(str, mapper, final)
  }

Each Routine function will be passed the current value and its index the
slice for its arguments.

*/
func Filter(data interface{}, routine Routine, callbacks ...Done) {
	var (
		routines []Routine
		results  []interface{}
	)

	d := reflect.ValueOf(data)

	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Interface()
		routines = append(routines, func(id int) Routine {
			return func(done Done, args ...interface{}) {
				done = func(original Done) Done {
					return func(err error, args ...interface{}) {
						if args[0] != false {
							results = append(results, v)
						}
						if id == (d.Len() - 1) {
							original(err, results...)
							return
						}
						original(err, args...)
					}
				}(done)

				routine(done, v, id)
			}
		}(i))
	}

	Waterfall(routines, callbacks...)
}

/*

FilterParallel allows you to filter out information from a slice in Parallel
mode.

You must call the Done function with false as its first argument if you do
not want the data to be present in the results. No other arguments will
affect the performance of this function.

If there is an error, any further results will be discarded but it will not
immediately exit. It will continue to run all of the other Routine functions
that were passed into it. This is because by the time the error is sent, the
goroutines have already been started. At this current time, there is no way
to cancel a sleep timer in Go.

For example, take a look at one of the tests for this function:
  func TestFilterStringParallel(t *testing.T) {
    str := []string{
      "test1",
      "test2",
      "test3",
      "test4",
      "test5",
    }

    expects := []string{
      "test1",
      "test2",
      "test4",
      "test5",
    }

    mapper := func(done async.Done, args ...interface{}) {
      Status("Hit string")
      Status("Args: %+v\n", args)
      if args[0] == "test3" {
        done(nil, false)
        return
      }
      done(nil, true)
    }

    final := func(err error, results ...interface{}) {
      Status("Hit string end")
      Status("Results: %+v\n", results)
      for i := 0; i < len(results); i++ {
        if results[i] != expects[i] {
          t.Errorf("Did not filter correctly.")
          break
        }
      }
    }

    async.FilterParallel(str, mapper, final)
  }

Each Routine function will be passed the current value and its index the
slice for its arguments.

The output of filtering in Parallel mode cannot be guaranteed to stay in the
same order, due to the fact that it may take longer to process some things
in your filter routine. If you need the data to stay in the order it is in,
use Filter instead to ensure it stays in order.

*/
func FilterParallel(data interface{}, routine Routine, callbacks ...Done) {
	var routines []Routine

	d := reflect.ValueOf(data)

	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Interface()
		routines = append(routines, func(id int) Routine {
			return func(done Done, args ...interface{}) {
				done = func(original Done) Done {
					return func(err error, args ...interface{}) {
						if args[0] != false {
							original(err, v)
							return
						}
						original(err)
					}
				}(done)

				routine(done, v, id)
			}
		}(i))
	}

	Parallel(routines, callbacks...)
}
