package async_test

import (
	"github.com/Southern/async"
	"testing"
)

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

func TestFilterBool(t *testing.T) {
	bools := []bool{
		true,
		false,
		false,
		true,
		false,
	}

	expects := []bool{
		true,
		false,
		true,
		false,
	}

	mapper := func(done async.Done, args ...interface{}) {
		Status("Hit bool")
		Status("Args: %+v\n", args)
		if args[1] == 2 {
			done(nil, false)
			return
		}
		done(nil, true)
	}

	final := func(err error, results ...interface{}) {
		Status("Hit bool end")
		Status("Results: %+v\n", results)
		for i := 0; i < len(results); i++ {
			if results[i] != expects[i] {
				t.Errorf("Did not filter correctly.")
				break
			}
		}
	}

	async.Filter(bools, mapper, final)
}

func TestFilterInt(t *testing.T) {
	bools := []int{
		1,
		2,
		3,
		4,
		5,
	}

	expects := []int{
		1,
		2,
		4,
		5,
	}

	mapper := func(done async.Done, args ...interface{}) {
		Status("Hit bool")
		Status("Args: %+v\n", args)
		if args[0] == 3 {
			done(nil, false)
			return
		}
		done(nil, true)
	}

	final := func(err error, results ...interface{}) {
		Status("Hit bool end")
		Status("Results: %+v\n", results)
		for i := 0; i < len(results); i++ {
			if results[i] != expects[i] {
				t.Errorf("Did not filter correctly.")
				break
			}
		}
	}

	async.Filter(bools, mapper, final)
}

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

func TestFilterBoolParallel(t *testing.T) {
	bools := []bool{
		true,
		false,
		false,
		true,
		false,
	}

	expects := []bool{
		true,
		false,
		true,
		false,
	}

	mapper := func(done async.Done, args ...interface{}) {
		Status("Hit bool")
		Status("Args: %+v\n", args)
		if args[1] == 2 {
			done(nil, false)
			return
		}
		done(nil, true)
	}

	final := func(err error, results ...interface{}) {
		Status("Hit bool end")
		Status("Results: %+v\n", results)
		for i := 0; i < len(results); i++ {
			if results[i] != expects[i] {
				t.Errorf("Did not filter correctly.")
				break
			}
		}
	}

	async.FilterParallel(bools, mapper, final)
}

func TestFilterIntParallel(t *testing.T) {
	bools := []int{
		1,
		2,
		3,
		4,
		5,
	}

	expects := []int{
		1,
		2,
		4,
		5,
	}

	mapper := func(done async.Done, args ...interface{}) {
		Status("Hit bool")
		Status("Args: %+v\n", args)
		if args[0] == 3 {
			done(nil, false)
			return
		}
		done(nil, true)
	}

	final := func(err error, results ...interface{}) {
		Status("Hit bool end")
		Status("Results: %+v\n", results)
		for i := 0; i < len(results); i++ {
			if results[i] != expects[i] {
				t.Errorf("Did not filter correctly.")
				break
			}
		}
	}

	async.FilterParallel(bools, mapper, final)
}
