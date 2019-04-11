package priority

import (
	"errors"
	"fmt"
	"reflect"
)

type Priority interface {
	ID() int64
	Name() string
	Deadline() Deadline
}

type Deadline struct {
	Duration           int
	Unit               string
	DeduceNonWorkHours bool
}

var priorities = make(map[int64]Priority, 0)

func RegisterPriority(new Priority) {
	vof := reflect.ValueOf(new)

	if reflect.Ptr != vof.Kind() {
		panic(errors.New("Registered document must be a pointer to struct"))
	}

	if _, exists := priorities[new.ID()]; exists {
		panic(fmt.Sprintf("priority %s already registered", new.Name()))
	}

	zero := reflect.New(reflect.TypeOf(vof.Elem().Interface()))

	priorities[new.ID()] = zero.Interface().(Priority)
}

func NewPriority(key int64) (Priority, bool) {
	data, exists := priorities[key]

	if !exists {
		return nil, false
	}

	new := reflect.New(reflect.TypeOf(data).Elem())

	return new.Interface().(Priority), true
}

func Priorities() map[int64]Priority {
	return priorities
}
