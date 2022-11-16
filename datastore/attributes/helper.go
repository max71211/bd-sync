package attributes

import (
	"strings"

	"github.com/pkg/errors"
)

type CheckTableFunc func(table string) bool

var ErrTableAttributeNotApproach = errors.New("attribute isn't approach")

func PrepareUpdateSets(check CheckTableFunc, attributes ...UpdateAttribute,
) (sets string, mapper map[string]interface{}, err error) {
	setParams := make([]string, 0, len(attributes))
	mapper = make(map[string]interface{})

	for _, attr := range attributes {
		if !check(attr.Table()) {
			return "", nil, errors.WithStack(ErrTableAttributeNotApproach)
		}
		setParams = append(setParams, attr.Query())
		mapper[attr.Key()] = attr.Value()
	}
	sets = strings.Join(setParams, ", ")
	return
}

func CheckTable(expectedTable string) CheckTableFunc {
	return func(table string) bool {
		return table == expectedTable
	}
}

func PrepareConditions(check CheckTableFunc, divider string, attributes ...ConditionAttribute) (sets string,
	mapper map[string]interface{}, err error) {
	setParams := make([]string, 0, len(attributes))
	mapper = make(map[string]interface{})

	for _, attr := range attributes {
		if !check(attr.Table()) {
			return "", nil, errors.WithStack(ErrTableAttributeNotApproach)
		}
		setParams = append(setParams, attr.Query())
		mapper[attr.Key()] = attr.Value()
	}
	sets = strings.Join(setParams, divider)
	return
}

func PrepareConditionsWithoutCheckTable(divider string, attributes ...ConditionAttribute) (sets string,
	mapper map[string]interface{}, err error) {
	setParams := make([]string, 0, len(attributes))
	mapper = make(map[string]interface{})

	for _, attr := range attributes {
		setParams = append(setParams, attr.Query())
		mapper[attr.Key()] = attr.Value()
	}
	sets = strings.Join(setParams, divider)
	return
}

func CombineMapper(mappers ...map[string]interface{}) map[string]interface{} {
	if len(mappers) == 0 {
		return make(map[string]interface{})
	}

	out := mappers[0]
	for _, mapper := range mappers {
		for key, value := range mapper {
			_, ok := out[key]
			if ok {
				continue
			}
			out[key] = value
			delete(mapper, key)
		}
	}
	return out
}
