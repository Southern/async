package filter

import (
  "git.aviuslabs.net/golang/async"
  "reflect"
)

func Filter(data interface{}, routine async.Routine, callbacks ...async.Done) {
  var (
    routines []async.Routine
    results  []interface{}
  )

  d := reflect.ValueOf(data)

  for i := 0; i < d.Len(); i++ {
    v := d.Index(i).Interface()
    routines = append(routines, func(done async.Done, args ...interface{}) {
      done = func(original async.Done) async.Done {
        return func(err error, args ...interface{}) {
          if args[0] != false {
            results = append(results, v)
          }
          if i == d.Len() {
            original(err, results...)
            return
          }
          original(err, args...)
        }
      }(done)

      routine(done, v)
    })
  }

  async.Waterfall(routines, callbacks...)
}
