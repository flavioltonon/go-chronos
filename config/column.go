package config

import (
	"errors"
	"fmt"
	"reflect"
)

type Column interface {
	ID() int64
	Name() string
	Project() int64
	StandardIssueState() string
}

var columns = make(map[int64]Column, 0)

func RegisterColumn(new Column) {
	vof := reflect.ValueOf(new)

	if reflect.Ptr != vof.Kind() {
		panic(errors.New("Registered column must be a pointer to struct"))
	}

	if _, exists := columns[new.ID()]; exists {
		panic(fmt.Sprintf("column %s already registered", new.Name()))
	}

	zero := reflect.New(reflect.TypeOf(vof.Elem().Interface()))

	columns[new.ID()] = zero.Interface().(Column)
}

func NewColumn(key int64) (Column, bool) {
	data, exists := columns[key]

	if !exists {
		return nil, false
	}

	new := reflect.New(reflect.TypeOf(data).Elem())

	return new.Interface().(Column), true
}

func Columns() map[int64]Column {
	return columns
}
