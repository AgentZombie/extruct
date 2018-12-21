package extruct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

type NotFoundError struct {
	Path string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("not found: %q", err.Path)
}

func (err NotFoundError) String() string {
	return err.Error()
}

func Extruct(v interface{}, path string) (interface{}, error) {
	return extruct(v, strings.Split(path, "/"), 0)
}

func extruct(v interface{}, path []string, offset int) (interface{}, error) {
	pathLen := len(path)
	if pathLen == 0 {
		return nil, ErrInvalidPath
	}
	if offset >= pathLen {
		return nil, errors.New("offset exceeds path")
	}
	if v == nil {
		if offset == pathLen-1 {
			return nil, nil
		}
		return nil, errors.New("nil in path at " + strings.Join(path[0:offset], "/"))
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	elemStr := path[offset]
	if len(elemStr) == 0 {
		return nil, ErrInvalidPath
	}
	field := val.FieldByName(elemStr)
	if field == (reflect.Value{}) {
		return nil, &NotFoundError{Path: strings.Join(path[0:offset+1], "/")}
	}
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		field = reflect.Indirect(field)
	}
	// terminal case
	if pathLen-1 == offset {
		return field.Interface(), nil
	}
	if field.Kind() != reflect.Slice {
		return extruct(field.Interface(), path, offset+1)
	}
	fLen := field.Len()
	vals := make([]interface{}, fLen)
	for i := 0; i < fLen; i++ {
		var err error
		vals[i], err = extruct(field.Index(i).Interface(), path, offset+1)
		if err != nil {
			return nil, err
		}
	}
	return vals, nil
}
