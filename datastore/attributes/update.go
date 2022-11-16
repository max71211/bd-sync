package attributes

import (
	"fmt"
)

// UpdateAttribute attribute for change field in database object
type UpdateAttribute interface {
	Key() string
	Value() interface{}
	Query() string
	Table() string
}

// NewTableUpdateAttribute simple update attribute, which help set 'SET' block in sql queries.
func NewTableUpdateAttribute(table, key string, value interface{}) UpdateAttribute {
	return &tableUpdateAttribute{
		key:   key,
		value: value,
		table: table,
	}
}

type tableUpdateAttribute struct {
	key   string
	value interface{}
	table string
}

func (attr tableUpdateAttribute) Key() string {
	return attr.key
}

func (attr tableUpdateAttribute) Value() interface{} {
	return attr.value
}

func (attr tableUpdateAttribute) Query() string {
	return fmt.Sprintf("%s = :%s", attr.Key(), attr.Key())
}

func (attr tableUpdateAttribute) Table() string {
	return attr.table
}

// NewTableQueryUpdateAttribute create update attributes, which value is come from query.
// Query must result one value
func NewTableQueryUpdateAttribute(table, key, query string) UpdateAttribute {
	return &queryTableUpdateAttributes{
		key:        key,
		valueQuery: query,
		table:      table,
	}
}

type queryTableUpdateAttributes struct {
	key        string
	table      string
	valueQuery string
}

func (attr queryTableUpdateAttributes) Key() string {
	return attr.key
}

func (attr queryTableUpdateAttributes) Value() interface{} {
	return nil
}

func (attr queryTableUpdateAttributes) Query() string {
	return fmt.Sprintf("%s = (%s)", attr.Key(), attr.valueQuery)
}

func (attr queryTableUpdateAttributes) Table() string {
	return attr.table
}
