package cli

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// mapInputSource implements FlagInputSourceExtension to return
// data from the map that is loaded.
type mapInputSource struct {
	file     string
	valueMap map[interface{}]interface{}
}

// Source returns the path of the source file
func (fsm *mapInputSource) Source() string {
	return fsm.file
}

// Bool returns an bool from the map otherwise returns false
func (fsm *mapInputSource) Bool(name string) (bool, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(bool)
		if !isType {
			return false, incorrectTypeForFlagError(name, "bool", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(bool)
		if !isType {
			return false, incorrectTypeForFlagError(name, "bool", nestedGenericValue)
		}
		return otherValue, nil
	}

	return false, nil
}

// Duration returns a duration from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Duration(name string) (time.Duration, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		return castDuration(name, otherGenericValue)
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		return castDuration(name, nestedGenericValue)
	}

	return 0, nil
}

func castDuration(name string, value interface{}) (time.Duration, error) {
	if otherValue, isType := value.(time.Duration); isType {
		return otherValue, nil
	}
	otherStringValue, isType := value.(string)
	parsedValue, err := time.ParseDuration(otherStringValue)
	if !isType || err != nil {
		return 0, incorrectTypeForFlagError(name, "duration", value)
	}
	return parsedValue, nil
}

// Float64 returns an float64 from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Float64(name string) (float64, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(float64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "float64", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(float64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "float64", nestedGenericValue)
		}
		return otherValue, nil
	}

	return 0, nil
}

// Float64Slice returns an []float64 from the map if it exists otherwise returns nil
func (fsm *mapInputSource) Float64Slice(name string) ([]float64, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if !exists {
		otherGenericValue, exists = nestedVal(name, fsm.valueMap)
		if !exists {
			return nil, nil
		}
	}

	otherValue, isType := otherGenericValue.([]interface{})
	if !isType {
		return nil, incorrectTypeForFlagError(name, "[]interface{}", otherGenericValue)
	}

	var float64Slice = make([]float64, 0, len(otherValue))
	for i, v := range otherValue {
		intValue, isType := v.(float64)

		if !isType {
			return nil, incorrectTypeForFlagError(fmt.Sprintf("%s[%v]", name, i), "float64", v)
		}

		float64Slice = append(float64Slice, intValue)
	}

	return float64Slice, nil
}

// Generic returns an cli.Generic from the map if it exists otherwise returns nil
func (fsm *mapInputSource) Generic(name string) (Generic, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(Generic)
		if !isType {
			return nil, incorrectTypeForFlagError(name, "cli.Generic", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(Generic)
		if !isType {
			return nil, incorrectTypeForFlagError(name, "cli.Generic", nestedGenericValue)
		}
		return otherValue, nil
	}

	return nil, nil
}

// Int returns an int from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Int(name string) (int, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(int)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "int", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(int)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "int", nestedGenericValue)
		}
		return otherValue, nil
	}

	return 0, nil
}

// IntSlice returns an []int from the map if it exists otherwise returns nil
func (fsm *mapInputSource) IntSlice(name string) ([]int, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if !exists {
		otherGenericValue, exists = nestedVal(name, fsm.valueMap)
		if !exists {
			return nil, nil
		}
	}

	otherValue, isType := otherGenericValue.([]interface{})
	if !isType {
		return nil, incorrectTypeForFlagError(name, "[]interface{}", otherGenericValue)
	}

	var intSlice = make([]int, 0, len(otherValue))
	for i, v := range otherValue {
		intValue, isType := v.(int)

		if !isType {
			return nil, incorrectTypeForFlagError(fmt.Sprintf("%s[%d]", name, i), "int", v)
		}

		intSlice = append(intSlice, intValue)
	}

	return intSlice, nil
}

// Int64 returns an int64 from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Int64(name string) (int64, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(int64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "int64", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(int64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "int64", nestedGenericValue)
		}
		return otherValue, nil
	}

	return 0, nil
}

// Int64Slice returns an []int64 from the map if it exists otherwise returns nil
func (fsm *mapInputSource) Int64Slice(name string) ([]int64, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if !exists {
		otherGenericValue, exists = nestedVal(name, fsm.valueMap)
		if !exists {
			return nil, nil
		}
	}

	otherValue, isType := otherGenericValue.([]interface{})
	if !isType {
		return nil, incorrectTypeForFlagError(name, "[]interface{}", otherGenericValue)
	}

	var int64Slice = make([]int64, 0, len(otherValue))
	for i, v := range otherValue {
		int64Value, isType := v.(int64)

		if !isType {
			return nil, incorrectTypeForFlagError(fmt.Sprintf("%s[%d]", name, i), "int64", v)
		}

		int64Slice = append(int64Slice, int64Value)
	}

	return int64Slice, nil
}

// String returns a string from the map if it exists otherwise returns an empty string
func (fsm *mapInputSource) String(name string) (string, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(string)
		if !isType {
			return "", incorrectTypeForFlagError(name, "string", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(string)
		if !isType {
			return "", incorrectTypeForFlagError(name, "string", nestedGenericValue)
		}
		return otherValue, nil
	}

	return "", nil
}

// StringSlice returns an []string from the map if it exists otherwise returns nil
func (fsm *mapInputSource) StringSlice(name string) ([]string, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if !exists {
		otherGenericValue, exists = nestedVal(name, fsm.valueMap)
		if !exists {
			return nil, nil
		}
	}

	otherValue, isType := otherGenericValue.([]interface{})
	if !isType {
		return nil, incorrectTypeForFlagError(name, "[]interface{}", otherGenericValue)
	}

	var stringSlice = make([]string, 0, len(otherValue))
	for i, v := range otherValue {
		stringValue, isType := v.(string)

		if !isType {
			return nil, incorrectTypeForFlagError(fmt.Sprintf("%s[%d]", name, i), "string", v)
		}

		stringSlice = append(stringSlice, stringValue)
	}

	return stringSlice, nil
}

// Uint returns an uint from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Uint(name string) (uint, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(uint)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "uint", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(uint)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "uint", nestedGenericValue)
		}
		return otherValue, nil
	}

	return 0, nil
}

// Uint64 returns an uint64 from the map if it exists otherwise returns 0
func (fsm *mapInputSource) Uint64(name string) (uint64, error) {
	otherGenericValue, exists := fsm.valueMap[name]
	if exists {
		otherValue, isType := otherGenericValue.(uint64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "uint64", otherGenericValue)
		}
		return otherValue, nil
	}
	nestedGenericValue, exists := nestedVal(name, fsm.valueMap)
	if exists {
		otherValue, isType := nestedGenericValue.(uint64)
		if !isType {
			return 0, incorrectTypeForFlagError(name, "uint64", nestedGenericValue)
		}
		return otherValue, nil
	}

	return 0, nil
}

// nestedVal checks if the name has '.' delimiters.
// If so, it tries to traverse the tree by the '.' delimited sections to find
// a nested value for the key.
func nestedVal(name string, tree map[interface{}]interface{}) (interface{}, bool) {
	if sections := strings.Split(name, "."); len(sections) > 1 {
		node := tree
		for _, section := range sections[:len(sections)-1] {
			child, ok := node[section]
			if !ok {
				return nil, false
			}
			ctype, ok := child.(map[interface{}]interface{})
			if !ok {
				return nil, false
			}
			node = ctype
		}
		if val, ok := node[sections[len(sections)-1]]; ok {
			return val, true
		}
	}
	return nil, false
}

func incorrectTypeForFlagError(name, expectedTypeName string, value interface{}) error {
	valueType := reflect.TypeOf(value)
	valueTypeName := ""
	if valueType != nil {
		valueTypeName = valueType.Name()
	}

	return fmt.Errorf("Mismatched type for flag '%s'. Expected '%s' but actual is '%s'", name, expectedTypeName, valueTypeName)
}
