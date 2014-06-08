package async

import (
  "reflect"
)

/*

  Filter out information from a slice in Waterfall mode.

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

      mapper := func(done Done, args ...interface{}) {
        Status("Hit string")
        Status("Args: %+v\n", args)
        if args[0] == "test3" {
          // We don't want this result in our return, so we send false back
          // as the first argument.
          done(nil, false)
          return
        }
        // We want anything else that we get, so we return true here.
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

      Filter(str, mapper, final)
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

  Filter out information from a slice in Parallel mode.

  You must call the Done function with false as its first argument if you do
  not want the data to be present in the results. No other arguments will
  affect the performance of this function. When calling the Done function,
  an error will cause the filtering to immediately exit.

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

      mapper := func(done Done, args ...interface{}) {
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

      FilterParallel(str, mapper, final)
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
            }
          }
        }(done)

        routine(done, v, id)
      }
    }(i))
  }

  Parallel(routines, callbacks...)
}
