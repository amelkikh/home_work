package hw09structvalidator

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrStringLen        = errors.New("invalid len")
	ErrValueNotInList   = errors.New("not exist")
	ErrNumberMin        = errors.New("value less than specified")
	ErrNumberMax        = errors.New("value more than specified")
	ErrRegexpNotMatched = errors.New("value not matched")

	ErrInvalidRegexp    = errors.New("invalid regexp")
	ErrParseValidator   = errors.New("invalid validator")
	ErrInvalidValueType = errors.New("invalid value type")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Field <= v[j].Field
	})
	for _, e := range v {
		b.WriteString(fmt.Sprintf("%q:%s\n", e.Field, e.Err.Error()))
	}
	return b.String()
}
