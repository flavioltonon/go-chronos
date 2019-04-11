package column

import "fmt"

type Column interface {
	ID() int64
	Name() string
}

var columns map[int64]Column

func NewColumn(column Column) {
	if _, exists := columns[column.ID()]; exists {
		panic(fmt.Sprintf("column %s already registered", column.Name()))
	}

	columns[column.ID()] = column
}
