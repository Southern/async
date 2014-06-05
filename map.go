package async

import (
  "reflect"
)

/*

  Manipulate data in a slice asynchronously.

  Each Routine will be called with the value and index of the current position
  in the slice. When calling the Done function, an error will cause the
  mapping to immediately exit. All other arguments are sent back as the
  replacement for the current value.

  For example, take a look at one of the tests for this function:
    func TestMapInt(t *testing.T) {
      ints := []int{1, 2, 3, 4, 5}

      expects := []int{2, 4, 6, 8, 10}

      mapper := func(done Done, args ...interface{}) {
        Status("Hit int")
        Status("Args: %+v\n", args)
        // We know we want an int. So we can cast this to an integer, multiply by
        // 2, and send it back as the result.
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

      Map(ints, mapper, final)
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
            if i == d.Len() {
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
