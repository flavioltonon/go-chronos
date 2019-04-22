package chronos

import (
	"errors"
	"fmt"
	"reflect"
)

type Project interface {
	ID() int64
	Name() string
}

var projects = make(map[int64]Project, 0)

func RegisterProject(new Project) {
	vof := reflect.ValueOf(new)

	if reflect.Ptr != vof.Kind() {
		panic(errors.New("Registered project must be a pointer to struct"))
	}

	if _, exists := projects[new.ID()]; exists {
		panic(fmt.Sprintf("project %s already registered", new.Name()))
	}

	zero := reflect.New(reflect.TypeOf(vof.Elem().Interface()))

	projects[new.ID()] = zero.Interface().(Project)
}

func NewProject(key int64) (Project, bool) {
	data, exists := projects[key]

	if !exists {
		return nil, false
	}

	new := reflect.New(reflect.TypeOf(data).Elem())

	return new.Interface().(Project), true
}

func Projects() map[int64]Project {
	return projects
}
