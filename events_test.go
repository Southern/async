package async_test

import (
	"fmt"
	"github.com/Southern/async"
	"testing"
	"time"
)

var events = make(async.Events)

func TestEventEmit(t *testing.T) {
	Status("Try nonexistant event")
	events.Emit("test")
}

func TestEventClearWithName(t *testing.T) {
	Status("Add event")
	events.On("test", func() {
		Status("Called test")
	}).On("test2", func() {
		Status("Called test2")
	})

	Status("Clear list")
	events.Clear("test")

	if events.Get("test") != nil {
		t.Errorf("Event was not properly removed.")
		return
	}
}

func TestEventClearNoName(t *testing.T) {
	Status("Add event")
	events.On("test", func() {
		Status("Called test")
	})

	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}
}

func TestEventOnWithoutArguments(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.On("test",
		func() {
			Status("Hit first callback")
		},
		func() {
			Status("Hit first callback")
		},
	)

	if events.Length("test") != 2 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test")

	if events.Length("test") != 2 {
		t.Errorf("One or more callbacks were removed")
		return
	}
}

func TestEventOnceWithoutArguments(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.Once("test",
		func() {
			Status("Hit first callback")
		},
		func() {
			Status("Hit second callback")
		},
	)

	if events.Length("test") != 2 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test")

	if events.Length("test") != 0 {
		t.Errorf("Not all callbacks were removed")
		return
	}
}

func TestEventOnWithArguments(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.On("test",
		func(msg string) {
			Status("Got message: %s", msg)
		},
	)

	if events.Length("test") != 1 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test", "blah")

	if events.Length("test") != 1 {
		t.Errorf("One or more callbacks were removed")
		return
	}
}

func TestEventOnceWithArguments(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.Once("test",
		func(msg string) {
			Status("Got message: %s", msg)
		},
	)

	if events.Length("test") != 1 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test", "blah")

	if events.Length("test") != 0 {
		t.Errorf("Not all callbacks were removed")
		return
	}
}

func TestEventMix(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.On("test",
		func() {
			Status("Hit first callback")
		},
	).Once("test",
		func() {
			Status("Hit second callback")
		},
	)

	if events.Length("test") != 2 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test")

	if events.Length("test") != 1 {
		t.Errorf("Callback should have been removed but wasn't")
		return
	}
}

func TestEventMixWithArguments(t *testing.T) {
	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.On("test",
		func(msg string) {
			Status("Hit first callback; msg: %s", msg)
		},
	).Once("test",
		func(msg string) {
			Status("Hit second callback; msg: %s", msg)
		},
	)

	if events.Length("test") != 2 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test", "Testing")

	if events.Length("test") != 1 {
		t.Errorf("Callback should have been removed but wasn't")
		return
	}
}

func TestEventErrorReturn(t *testing.T) {
	var _err error

	Status("Clear list")
	events.Clear()

	if len(events) > 0 {
		t.Errorf("List wasn't cleared")
		return
	}

	Status("Add event")
	events.On("test",
		func() error {
			return fmt.Errorf("Testing")
		},
	).On("error",
		func(err error) {
			Status("Got error: %s", err)
			_err = err
		},
	)

	if events.Length("test") != 1 {
		t.Errorf("Not all callbacks were added")
		return
	}

	Status("Emitting event")
	events.Emit("test")

	if events.Length("test") != 1 {
		t.Errorf("One or more callbacks were removed")
		return
	}

	Status("Giving the error time to trigger")
	time.Sleep(time.Second)

	if _err == nil {
		t.Errorf("Expected an error")
		return
	}
}

type MyStruct struct {
	async.Events
}

func TestEventInheritanceWithoutArguments(t *testing.T) {
	Status("Creating struct")
	mystruct := MyStruct{make(async.Events)}

	Status("Adding event")
	mystruct.On("test", func() {
		Status("Hit callback")
	})

	Status("Emitting event")
	mystruct.Emit("test")
}

func TestEventInheritanceWithArguments(t *testing.T) {
	Status("Creating struct")
	mystruct := MyStruct{make(async.Events)}

	Status("Adding event")
	mystruct.On("test", func(msg string) {
		Status("Hit callback with message: %s", msg)
	})

	Status("Emitting event")
	mystruct.Emit("test", "Testing")
}

func TestEventInheritanceErrorReturn(t *testing.T) {
	var _err error

	Status("Creating struct")
	mystruct := MyStruct{make(async.Events)}

	Status("Adding event")
	mystruct.On("error", func(err error) {
		_err = err
		Status("Got error: %s", err)
	}).On("test", func(msg string) error {
		return fmt.Errorf("Testing")
	})

	Status("Emitting event")
	mystruct.Emit("test", "Testing")

	Status("Giving the error time to trigger")
	time.Sleep(time.Second)

	if _err == nil {
		t.Errorf("Expected an error")
		return
	}
}
