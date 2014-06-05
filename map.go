package async

import (
  "reflect"
)

func Map(data interface{}, routine Routine, callbacks ...Done) {
  var (
    routines []Routine
    results  []interface{}
  )

  d := reflect.ValueOf(data)

  for i := 0; i < d.Len(); i++ {
    v := d.Index(i).Interface()
    routines = append(routines, func(done Done, args ...interface{}) {
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

      routine(done, v)
    })
  }

  Waterfall(routines, callbacks...)
}
