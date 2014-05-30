/*
  Asynchronous function stack that simplifies asynchronous Go code for the
  masses.

    list := async.New()
    list.Add([]async.Routine{
      func(done async.Done, args ...interface{}) {
        fmt.Printf("Arguments: %+v\n", args)
        done(nil, "result 1", "result 2", "result 3")
      },

      func(done async.Done, args ...interface{}) {
        fmt.Printf("Arguments: %+v\n", args)
        done(nil, "result 3", "result 4", "result 5")
      },
    })

    err, list := list.Run(func(err error, results ...interface{}) {
      fmt.Printf("Results: %+v\n", results)
    })

    if err != nil {
      fmt.Errorf("Got error: %s\n", error)
      return
    }
    println("Running async functions")
*/
package async

import (
  "fmt"
  "sync"
)

/*
  Done functions are called once an Routine function is done.

  It should accept an error as its first argument, and all other arguments
  should be handled as results from the Routine function.
*/
type Done func(error, ...interface{})

/*
  Routine functions are the main functions that will be called when
  running through the Routines array.

  It should accept an Done function as its first argument, and all
  other arguments will be passed into the Routine function.

  Every Routine function must call its Done function. The Done
  function will automatically return if there has been an error detected from
  a previous Routine function.

    list := async.New()
    list.Add(func(done async.Done, args ...interface{}) {
      fmt.Printf("Arguments: %+v\n", args)
      done(nil, "result 1", "result 2", "result 3")
    })
    list.Run(func(err error, results ...interface{}) {
      fmt.Printf("Results: %+v\n", results)
    })

  This allows us to add multiple callbacks easily.

    list.Add([]async.Routine{
      func(done async.Done, args ...interface{}) {
        // Your first callback
      },

      func(done async.Done, args ...interface{}) {
        // Your second callback
      },
    })
*/
type Routine func(Done, ...interface{})

/*
  Error that is used when there are no Routine functions to call.

  This error is only used when the list is executed with an empty Routines
  list.

  For example, Run and RunParallel would return this error if there
  are no Routine functions when they are called.
*/
var NoRoutines = fmt.Errorf("There are no routines to run.")

/*
  Asynchronous list that contains the sync.WaitGroup and multiple Routine
  functions.
*/
type List struct {
  Wait     sync.WaitGroup
  Routines []Routine
}

/*
  Create a new asynchronous list that contains the sync.WaitGroup and
  multiple Routine functions.
*/
func New() *List {
  return &List{
    Wait: sync.WaitGroup{},
  }
}

/*
  Add one or more Routine functions to the current list.
*/
func (list *List) Add(routines ...Routine) (error, *List) {
  list.Routines = append(list.Routines, routines...)
  return nil, list
}

/*
  Clear all Routine functions from the current list.
*/
func (list *List) Clear() *List {
  // Make sure the routines are empty
  list.Routines = list.Routines[:0]

  // Return an entirely new list so we don't have to cancel waits and such
  return New()
}

/*
  Used in Run() as the Done callback to call the next Routine in
  the stack.
*/
func (list *List) nextWaterfall(callbacks ...Done) Done {
  // final callback to trigger all callbacks that were set to happen after
  // the waterfall event
  final := func(err error, args ...interface{}) {
    for i := 0; i < len(callbacks); i++ {
      callbacks[i](err, args...)
    }
  }

  return func(err error, args ...interface{}) {
    list.Wait.Done()

    if err != nil {
      // Send the error off to final
      final(err, list)

      // Cancel all other waits
      for i := 0; i < (len(list.Routines) - 1); i++ {
        list.Wait.Done()
      }

      // Exit out of this function so we don't try any more routines
      return
    }

    if len(list.Routines) > 1 {
      // Fire off the next routine
      go func() {
        // Run the actual routine
        list.Routines[0](list.nextWaterfall(callbacks...), args...)

        // Shift the routine array
        list.Routines = list.Routines[1:]

        // Make our wait group actually triggers the wait
        list.Wait.Wait()
      }()

      // Exit out of this shit so it doesn't trigger the final callbacks
      return
    }

    // Send the arguments over to final
    final(nil, args...)
  }
}

/*
  Run all Routine functions in Waterfall mode.
*/
func (list *List) Run(callbacks ...Done) (error, *List) {
  if len(list.Routines) == 0 {
    for i := 0; i < len(callbacks); i++ {
      callbacks[i](NoRoutines, nil)
    }
    return NoRoutines, list
  }

  // Add all of our routines to the current wait group.
  list.Wait.Add(len(list.Routines))

  // Start off the first routine
  go func() {
    // Run the actual routine
    list.Routines[0](list.nextWaterfall(callbacks...), nil)

    // Shift the routine array
    list.Routines = list.Routines[1:]
  }()

  // Make our wait group actually triggers the wait
  list.Wait.Wait()

  return nil, list
}

/*
  Run all Routine functions in Parallel mode.
*/
func (list *List) RunParallel(callbacks ...Done) (error, *List) {
  var (
    // Interface array to hold results from async functions
    results = make([]interface{}, 0)

    // Channel to receive and send arguments from async functions
    result = make(chan interface{})

    // Did we get an error somewhere? This is used to exit when an error is
    // found
    _error = false

    // final callback to trigger all callbacks that were set to happen after
    // the waterfall event
    final = func(err error, args ...interface{}) {
      for i := 0; i < len(callbacks); i++ {
        callbacks[i](err, args...)
      }
    }
  )

  if len(list.Routines) == 0 {
    final(NoRoutines, nil)
    return NoRoutines, list
  }

  // Add all of our routines to the current wait group.
  list.Wait.Add(len(list.Routines))

  // Run all routines parallel
  for i := 0; i < len(list.Routines); i++ {
    go func(routine Routine) {
      routine(func(err error, args ...interface{}) {
        if err != nil {
          result <- err
          return
        }

        result <- args
      }, results...)

      list.Wait.Done()
    }(list.Routines[i])
  }

  // Clear out the routines so they can't be ran again
  list.Routines = list.Routines[:0]

  // Handle the results from the result channel
  go func() {
    for {
      r := <-result

      switch r.(type) {
      case error:
        _error = true
        final(r.(error), nil)

      case []interface{}:
        results = append(results, r.([]interface{})...)
      }
    }
  }()

  // Make our wait group actually triggers the wait
  list.Wait.Wait()

  // Failsafe to prevent callback from being fired again if an error occurred
  if !_error {
    final(nil, results...)
  }

  return nil, list
}

/*
  Run each Routine function in order, pass the arguments from the
  previous Routine to the next, and wait for the Done function
  to be called before continuing to the next Routine function.

  The callbacks will be called with the arguments from the last Routine
  that is called. This means that you will have to join the arguments together
  on your own.
*/
func Waterfall(routines []Routine, callbacks ...Done) {
  // Create a new list
  list := New()

  // Add all of the routines
  list.Add(routines...)

  list.Run(callbacks...)
}

/*
  Run each Routine function in order, pass the arguments from the
  previous Routine to the next, but do not wait for the previous
  function to call its Done function.

  Unlike Waterfall, which doesn't combine the arguments returned from its
  Routine functions, Parallel does combine the arguments returned from
  its Routine functions. This means that the callbacks will have the
  results from every Routine function when they are called.
*/
func Parallel(routines []Routine, callbacks ...Done) {
  // Create a new list
  list := New()

  // Add all of the routines
  list.Add(routines...)

  list.RunParallel(callbacks...)
}
