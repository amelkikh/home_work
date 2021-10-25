package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator interface {
	Validate(interface{}) (bool, error)
}

var TagName = "validate"

type validator struct {
	errors ValidationErrors
}

func (v *validator) validateString(field string, tags []string, str string) error {
	switch tags[0] {
	case "len":
		lenValue, err := strconv.Atoi(tags[1])
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
		}
		if len(str) != lenValue {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrStringLen})
		}
	case "in":
		in := strings.Split(tags[1], ",")
		m := make(map[string]struct{})
		for _, vv := range in {
			m[vv] = struct{}{}
		}
		if _, ok := m[str]; !ok {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrValueNotInList})
		}
	case "regexp":
		re, err := regexp.Compile(tags[1])
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrInvalidRegexp)
		}
		if !re.MatchString(str) {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrRegexpNotMatched})
		}
	}
	return nil
}

func (v *validator) validateSlice(field string, tags []string, val reflect.Value) error {
	if val.Len() > 0 {
		sValType := val.Index(0).Kind()
		switch sValType {
		case reflect.String:
			for j := 0; j < val.Len(); j++ {
				err := v.validateString(field, tags, val.Index(j).String())
				if err != nil {
					return err
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for j := 0; j < val.Len(); j++ {
				err := v.validateInt(field, tags, val.Index(j).Int())
				if err != nil {
					return err
				}
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for j := 0; j < val.Len(); j++ {
				err := v.validateUInt(field, tags, val.Index(j).Uint())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *validator) validateUInt(field string, tags []string, i uint64) error {
	switch tags[0] {
	case "min":
		val, err := strconv.ParseUint(tags[1], 10, 64)
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
		}
		if i < val {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrNumberMin})
		}
	case "max":
		val, err := strconv.ParseUint(tags[1], 10, 64)
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
		}
		if i > val {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrNumberMax})
		}
	case "in":
		in := strings.Split(tags[1], ",")
		m := make(map[uint64]struct{})
		for _, vv := range in {
			val, err := strconv.ParseUint(vv, 10, 64)
			if err != nil {
				return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
			}
			m[val] = struct{}{}
		}
		if _, ok := m[i]; !ok {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrValueNotInList})
		}
	}
	return nil
}

func (v *validator) validateInt(field string, tags []string, i int64) error {
	switch tags[0] {
	case "min":
		val, err := strconv.ParseInt(tags[1], 10, 64)
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
		}
		if i < val {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrNumberMin})
		}
	case "max":
		val, err := strconv.ParseInt(tags[1], 10, 64)
		if err != nil {
			return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
		}
		if i > val {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrNumberMax})
		}
	case "in":
		in := strings.Split(tags[1], ",")
		m := make(map[int64]struct{})
		for _, vv := range in {
			val, err := strconv.ParseInt(vv, 10, 64)
			if err != nil {
				return fmt.Errorf("value %q: %w", tags[1], ErrParseValidator)
			}
			m[val] = struct{}{}
		}
		if _, ok := m[i]; !ok {
			v.errors = append(v.errors, ValidationError{Field: field, Err: ErrValueNotInList})
		}
	}
	return nil
}

func Validate(v interface{}) error {
	val := &validator{}
	return val.Validate(v)
}

func (v *validator) Validate(vv interface{}) error {
	refValue := reflect.ValueOf(vv)
	refType := reflect.TypeOf(vv)

	if refType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: value %s", ErrInvalidValueType, refType.Kind().String())
	}
	refValueType := refValue.Type()
	for i := 0; i < refValue.NumField(); i++ {
		refStructField := refValueType.Field(i)
		fName := refStructField.Name
		tagStr, ok := refValueType.Field(i).Tag.Lookup(TagName)
		if !ok {
			continue
		}
		if strings.Compare(tagStr, "nested") == 0 {
			if refValue.Field(i).Elem().CanInterface() {
				err := v.Validate(refValue.Field(i).Elem().Interface())
				if err != nil {
					return err
				}
				continue
			}
		}

		validators := strings.Split(tagStr, "|")
		for _, vldtr := range validators {
			tags := strings.Split(vldtr, ":")
			if len(tags) != 2 || len(tags[0]) == 0 || len(tags[1]) == 0 {
				return ErrParseValidator
			}

			switch refValue.Field(i).Kind() {
			case reflect.Slice:
				if err := v.validateSlice(fName, tags, refValue.Field(i)); err != nil {
					return err
				}
			case reflect.String:
				if err := v.validateString(fName, tags, refValue.Field(i).String()); err != nil {
					return err
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if err := v.validateInt(fName, tags, refValue.Field(i).Int()); err != nil {
					return err
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if err := v.validateUInt(fName, tags, refValue.Field(i).Uint()); err != nil {
					return err
				}
			}
		}
	}

	if len(v.errors) > 0 {
		return v.errors
	}
	return nil
}
