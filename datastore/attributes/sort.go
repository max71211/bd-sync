package attributes

import (
	"fmt"
	"strings"
)

const querySort = "ORDER BY"

// SortAttribute setting sorting
type SortAttribute struct {
	Table      string
	StrictFunc func(string) string
	Key        string
	Direct     bool
}

// WithStrict is a strict table name formatter that wraps the table name in quotation marks
func WithStrict(unit string) string {
	return fmt.Sprintf(`"%s"`, unit)
}

// WithoutStrict is a non-strict formatter for the table name,
// which is the default for formatting that returns the incoming element
func WithoutStrict(unit string) string {
	return unit
}

// String implement fmt.Stringer
func (attr SortAttribute) String() string {
	if attr.Key == "" {
		return ""
	}

	return fmt.Sprintf(fmt.Sprintf("%s %s", querySort, attr.Order()))
}

// Order format order attribute
func (attr SortAttribute) Order() string {
	if attr.StrictFunc == nil {
		attr.StrictFunc = WithoutStrict
	}

	direction := "ASC"
	if !attr.Direct {
		direction = "DESC"
	}

	column := attr.Key
	if attr.Table != "" {
		column = fmt.Sprintf(`%s.%s`, attr.StrictFunc(attr.Table), attr.Key)
	}

	return fmt.Sprintf("%s %s", column, direction)
}

// MultiSortAttribute multiple sort attributes
type MultiSortAttribute []SortAttribute

// String implement fmt.Stringer
func (attr MultiSortAttribute) String() string {
	if len(attr) == 0 {
		return ""
	}

	exprs := make([]string, 0, len(attr))
	for _, item := range attr {
		if item.Key == "" {
			continue
		}
		exprs = append(exprs, item.Order())
	}

	return fmt.Sprintf("%s %s", querySort, strings.Join(exprs, ","))
}
