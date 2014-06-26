/*

Package async is a package to provide simplistic asynchronous routines for the
masses.

*/
package async

/*

Done types are used for shorthand definitions of the functions that are
passed into each Routine to show that the Routine has completed.

An example a Done function would be:
  func ImDone(err error, args ...interface{}) {
    if err != nil {
      // Handle the error your Routine returned.
      return
    }

    // There wasn't an error returned your Routine! Do what you want with
    // the args.
  }

*/
type Done func(error, ...interface{})

/*

Routine types are used for shorthand definitions of the functions that are
actually ran when calling Parallel, Waterfall, etc.

An example of a Routine function would be:
  func MyRoutine(done async.Done, args ...interface{}) {
    // Do something in your routine and then call its done function.
    done(nil, "arg1", "arg2", "arg3")
  }

*/
type Routine func(Done, ...interface{})
