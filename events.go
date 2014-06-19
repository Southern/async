package async

import (
  "reflect"
)

/*

  Event is a map of functions to use for a single event. The integer is the
  frequency of calls that the function should make.

*/
type Event map[reflect.Value]int

/*

  List is used for containing everything related to the events. It stores the
  name of the event, functions for the event, and number of times the function
  should be called when the event is triggered.

  All you need to do to create an event list is:
    events := make(async.Events)

  All event commands that return Events can be chained together. For example:
    events.On("myevent", func() {
      println("Called myevent")
    }).On("myevent2", func(msg string) {
      fmt.Printf("Called myevent2 with message: %s\n", msg)
    }).Emit("myevent").Emit("myevent2", "Testing")

  You can return an error from a function and it will be emitted as an error
  event. For example:
    events.On("error", func(err error) {
      fmt.Printf("Error: %s", err)
    }).On("myevent", func() error {
      return fmt.Errorf("Some error message")
    }).Emit("myevent")


  It's also easily inheritable by other structures. For example:
    type MyStruct struct {
      Events
    }

    m := MyStruct{make(async.Events)}
    m.On("myevent", func() {
      println("Called myevent")
    }).Emit("myevent")

*/
type Events map[string]Event

/*

  Clear all events out of the event list. You can supply optional names for
  the events to be cleared.

  For instance:
    events.On("test", func() {}).On("test2", func() {}).Emit("test").Clear("test")

  Returns the list of events for chaining commands.

*/
func (e Events) Clear(name ...string) Events {
  if name != nil {
    for i := 0; i < len(name); i++ {
      delete(e, name[i])
    }
    return e
  }

  for key := range e {
    delete(e, key)
  }
  return e
}

/*

  Emit an event. Arguments are optional. Each event will be ran as a Series.

  For example:
    events := make(async.Events)
    events.On("myevent", func() {
      println("Emitted myevent")
    })
    events.Emit("myevent")

  With arguments:
    events := make(async.Events)
    events.On("myevent", func(msg string) {
      fmt.Printf("Message: %s\n", msg)
    })
    events.Emit("myevent", "Testing")

  Returns the list of events for chaining commands.

*/
func (e Events) Emit(name string, args ...interface{}) Events {
  var (
    routines = make([]Routine, 0)
    values   = make([]reflect.Value, 0)
  )

  // If we don't have any events with this name, simply return.
  if e.Get(name) == nil {
    return e
  }

  // Reflect all of our arguments for the reflect.Value.Call
  for i := 0; i < len(args); i++ {
    values = append(values, reflect.ValueOf(args[i]))
  }

  for fn, freq := range e[name] {
    // Decrease frequency
    if freq > 0 {
      freq--
    }

    // If the frequency is down to 0, remove the callback from the event
    // so that it isn't triggered again.
    if freq == 0 {
      delete(e[name], fn)
    } else {
      e[name][fn] = freq
    }

    // Delete the entire event if all callbacks have been triggered
    if len(e[name]) == 0 {
      delete(e, name)
    }

    // Create the routines to pass into Series
    routines = append(routines,
      func(e Events, fn reflect.Value, values []reflect.Value) Routine {
        return func(done Done, args ...interface{}) {
          values := fn.Call(values)
          for i := 0; i < len(values); i++ {
            v := values[i].Interface()
            switch v.(type) {
            case error:
              done(v.(error))
              return
            }
          }
          done(nil)
        }
      }(e, fn, values),
    )
  }

  // Run all of the events in Series
  Series(routines, func(err error, args ...interface{}) {
    // Only emit the error event if an error was detected. Nothing else needs
    // to be done here.
    if err != nil {
      e.Emit("error", err)
    }
  })

  return e
}

/*

  Get Event map of functions and frequencies for the named event. This is just
  a convenience function. This data could also be accessed by the normal
  mapping methods.

  For instance:
    fmt.Printf("Events for myevent: %+v\n", e["myevent"])

*/
func (e Events) Get(name string) Event {
  return e[name]
}

/*

  Get length of the Event map of functions and frequencies for the named
  event. This is just a convenience function. This data could also be accessed
  by the normal mapping methods.

  For instance:
    fmt.Printf("Length: %d", len(e["myevent"]))

*/
func (e Events) Length(name string) int {
  return len(e.Get(name))
}

/*

  Add an event to be called forever.

  This is equal to calling Times with -1 as the number of times to run the
  event. More documentation can be found on the Times function.

  Returns the list of events for chaining commands.

*/
func (e Events) On(name string, callbacks ...interface{}) Events {
  return e.Times(name, -1, callbacks...)
}

/*

  Add an event to be called once.

  This is equal to calling Times with 1 as the number of times to run the
  event. More documentation can be found on the Times function.

  Returns the list of events for chaining commands.

*/
func (e Events) Once(name string, callbacks ...interface{}) Events {
  return e.Times(name, 1, callbacks...)
}

/*

  Add an event to be called a number of times. If the number of times for the
  function to be called is -1, it will be called until the list is cleared.

  Returns the list of events for chaining commands.

*/
func (e Events) Times(name string, times int, callbacks ...interface{}) Events {
  // Check to see if the event already exists. If not, create its map.
  if e[name] == nil {
    e[name] = make(Event)
  }

  for i := 0; i < len(callbacks); i++ {
    // Reflect the function so that we don't have to add function restraints.
    fn := reflect.ValueOf(callbacks[i])

    // Set the number of times that the event should run.
    e[name][fn] = times
  }

  return e
}
