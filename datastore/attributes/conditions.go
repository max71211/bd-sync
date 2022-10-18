package attributes

import (
	"fmt"
	"strings"
	"time"
)

type ConditionAttribute interface {
	Key() string
	Value() interface{}
	Query() string
	Table() string
}

func NewCondition(table, key string, cond string, value interface{}) *Condition {
	return &Condition{
		key:       key,
		value:     value,
		condition: cond,
		table:     table,
	}
}

func NewStrictCondition(table, key string, value interface{}) *Condition {
	return NewCondition(table, key, "=", value)
}

func NewLessCondition(table, key string, value interface{}) *Condition {
	return NewCondition(table, key, "<", value)
}

func NewStrictLessCondition(table, key string, value interface{}) *Condition {
	return NewCondition(table, key, "<=", value)
}

func NewMoreCondition(table, key string, value interface{}) *Condition {
	return NewCondition(table, key, ">", value)
}

func NewStrictMoreCondition(table, key string, value interface{}) *Condition {
	return NewCondition(table, key, ">=", value)
}

type Condition struct {
	key       string
	value     interface{}
	condition string
	table     string

	ValueKey   string
	StrictFunc func(string) string
}

func (attr *Condition) Key() string {
	if attr.ValueKey == "" {
		key := strings.Replace(attr.key, "\"", "", -1)
		attr.ValueKey = fmt.Sprintf("%s%d", key, time.Now().UnixNano())
	}
	return attr.ValueKey
}

func (attr *Condition) Value() interface{} {
	return attr.value
}

func (attr *Condition) Query() string {
	if attr.StrictFunc == nil {
		attr.StrictFunc = WithoutStrict
	}

	key := attr.key
	if attr.table != "" {
		key = fmt.Sprintf("%s.%s", attr.StrictFunc(attr.table), attr.key)
	}
	return fmt.Sprintf("%s %s :%s", key, attr.condition, attr.Key())
}

func (attr *Condition) Table() string {
	return attr.table
}

func NewSliceCondition(table string, key string, values ...interface{}) ConditionAttribute {
	return &sliceCondition{
		table: table,
		query: fmt.Sprintf("%s IN (:%s)", key, key),
		key:   key,
		value: values,
	}
}

func NewNotSliceCondition(table string, key string, values ...interface{}) ConditionAttribute {
	return &sliceCondition{
		table: table,
		query: fmt.Sprintf("%s NOT IN (:%s)", key, key),
		key:   key,
		value: values,
	}
}

type sliceCondition struct {
	table string
	query string
	key   string
	value interface{}
}

func (attr *sliceCondition) Key() string {
	return attr.key
}

func (attr *sliceCondition) Value() interface{} {
	return attr.value
}

func (attr *sliceCondition) Query() string {
	return attr.query
}

func (attr *sliceCondition) Table() string {
	return attr.table
}

func NewNullableChecker(table, key string, isSet bool) *NullableChecker {
	return &NullableChecker{
		key:   key,
		table: table,
		isSet: isSet,
	}
}

type NullableChecker struct {
	key   string
	table string
	isSet bool

	ValueKey string
}

func (attr *NullableChecker) Key() string {
	if attr.ValueKey == "" {
		attr.ValueKey = fmt.Sprintf("%s%d", attr.key, time.Now().UnixNano())
	}
	return attr.ValueKey
}

func (attr *NullableChecker) Value() interface{} {
	return nil
}

func (attr *NullableChecker) Query() string {
	if attr.isSet {
		return fmt.Sprintf("%s IS NULL", attr.key)
	}
	return fmt.Sprintf("%s IS NOT NULL", attr.key)
}

func (attr *NullableChecker) Table() string {
	return attr.table
}
