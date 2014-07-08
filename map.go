package async

import (
	"reflect"
)

/*

Map allows you to manipulate data in a slice in Waterfall mode.

Each Routine will be called with the value and index of the current position
in the slice. When calling the Done function, an error will cause the
mapping to immediately exit. All other arguments are sent back as the
replacement for the current value.

For example, take a look at one of the tests for this function:
  func TestMapInt(t *testing.T) {
    ints := []int{1, 2, 3, 4, 5}

    expects := []int{2, 4, 6, 8, 10}

    mapper := func(done async.Done, args ...interface{}) {
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

    async.Map(ints, mapper, final)
  }

*/
func Map(data interface{}, routine Routine, callbacks ...Done) {
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
						results = append(results, args...)
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

MapParallel allows you to manipulate data in a slice in Parallel mode.

Each Routine will be called with the value and index of the current position
in the slice. When calling the Done function, arguments are sent
back as the replacement for the current value.

If there is an error, any further results will be discarded but it will not
immediately exit. It will continue to run all of the other Routine functions
that were passed into it. This is because by the time the error is sent, the
goroutines have already been started. At this current time, there is no way
to cancel a sleep timer in Go.

For example, take a look at one of the tests for this function:
  func TestMapStringParallel(t *testing.T) {
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

    mapper := func(done async.Done, args ...interface{}) {
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

    async.MapParallel(str, mapper, final)
  }

The output of mapping in Parallel mode cannot be guaranteed to stay in the
same order, due to the fact that it may take longer to process some things
in your map routine. If you need the data to stay in the order it is in, use
Map instead to ensure it stays in order.

*/
func MapParallel(data interface{}, routine Routine, callbacks ...Done) {
	var routines []Routine

	d := reflect.ValueOf(data)

	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Interface()
		routines = append(routines, func(id int) Routine {
			return func(done Done, args ...interface{}) {
				routine(done, v, id)
			}
		}(i))
	}

	Parallel(routines, callbacks...)
}
