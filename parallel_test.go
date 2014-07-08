package async_test

import (
	"fmt"
	"github.com/Southern/async"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	Status("Calling parallel")
	async.Parallel([]async.Routine{
		func(done async.Done, args ...interface{}) {
			Status("First parallel function")
			Status("Called with arguments: %+v", args)
			done(nil, "arg1", "arg2", "arg3")
		},

		func(done async.Done, args ...interface{}) {
			Status("Second parallel function")
			Status("Called with arguments: %+v", args)
			time.Sleep(time.Second)
			done(nil, "arg4", "arg5", "arg6")
		},

		func(done async.Done, args ...interface{}) {
			Status("Third parallel function")
			Status("Called with arguments: %+v", args)
			done(nil, "arg7", "arg8", "arg9")
		},
	}, func(err error, results ...interface{}) {
		if err != nil {
			t.Errorf("Parallel threw an unexpected error: %+v", err)
			return
		}

		Status("Parallel completed with results: %+v", results)
	})
}

func TestParallelError(t *testing.T) {
	Status("Calling Parallel")
	async.Parallel([]async.Routine{
		func(done async.Done, args ...interface{}) {
			Status("First parallel function")
			Status("Called with arguments: %+v", args)
			time.Sleep(time.Second)
			done(nil, "arg1", "arg2", "arg3")
		},

		func(done async.Done, args ...interface{}) {
			Status("Second parallel function")
			Status("Called with arguments: %+v", args)
			done(fmt.Errorf("Test error"))
		},

		func(done async.Done, args ...interface{}) {
			Status("Third parallel function")
			Status("Called with arguments: %+v", args)
			done(nil, "arg4", "arg5", "arg6")
		},
	}, func(err error, results ...interface{}) {
		if err != nil {
			Status("Parallel exited with error: %+v", err)
			return
		}

		t.Errorf("Parallel did not throw an error as expected")
	})
}
