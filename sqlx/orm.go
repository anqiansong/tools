package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	errInvalidPointer = errors.New("invalid pointer")
	errNotSettable    = errors.New("can not settable")

	// ErrNoRows is an alias of sql.ErrNoRows
	ErrNoRows = sql.ErrNoRows
)

// UnmarshalRow accepts an interface to scan in, there is one row could be scan even though
// there are more than one rows, an ErrNoRows will be returned if have no rows.
func UnmarshalRow(rows *sql.Rows, v interface{}) error {
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return ErrNoRows
	}

	if err := must(v); err != nil {
		return err
	}

	t := reflect.TypeOf(v)
	it := indirect(t)
	switch it.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8,
		reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32,
		reflect.Float64, reflect.String:
		return scanBasicRow(rows, v)
	case reflect.Struct:
		return scanStructRow(rows, v)
	default:
		return errors.New("unsupported type")
	}
}

// UnmarshalRows accepts an interface which type must be ptr-slice
func UnmarshalRows(rows *sql.Rows, v interface{}) error {
	if err := must(v); err != nil {
		return err
	}

	// ptr *[]int
	ptr := reflect.TypeOf(v)
	ptrv := reflect.ValueOf(v)
	// slice []int
	slice := ptr.Elem()
	slicev := reflect.Indirect(ptrv)
	if !slicev.CanSet() {
		return errNotSettable
	}

	item := slice.Elem()
	itemBaseType := indirect(item)
	switch itemBaseType.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8,
		reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32,
		reflect.Float64, reflect.String:
		for rows.Next() {
			value := reflect.New(itemBaseType)
			err := rows.Scan(value.Interface())
			if err != nil {
				return err
			}

			if item.Kind() == reflect.Ptr {
				slicev.Set(reflect.Append(slicev, value))
			} else {
				slicev.Set(reflect.Append(slicev, reflect.Indirect(value)))
			}
		}
	case reflect.Struct:
		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			value := reflect.New(itemBaseType)
			list, err := convertStructFieldsIntoInterfaceSlice(value, columns)
			if err != nil {
				return err
			}

			err = rows.Scan(list...)
			if err != nil {
				return err
			}

			if item.Kind() == reflect.Ptr {
				slicev.Set(reflect.Append(slicev, value))
			} else {
				slicev.Set(reflect.Append(slicev, reflect.Indirect(value)))
			}
		}
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func scanStructRow(rows *sql.Rows, v interface{}) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	value := reflect.ValueOf(v)
	list, err := convertStructFieldsIntoInterfaceSlice(value, columns)
	if err != nil {
		return err
	}

	return rows.Scan(list...)
}

func convertStructFieldsIntoInterfaceSlice(v reflect.Value, columns []string) ([]interface{}, error) {
	fields, err := getFields(v)
	if err != nil {
		return nil, err
	}

	if len(columns) < len(fields) {
		return nil, fmt.Errorf("expected column num %d, but found %d", len(columns), len(fields))
	}

	list := make([]interface{}, 0)
	for _, column := range columns {
		if v, ok := fields[column]; ok {
			list = append(list, v.Interface())
		} else {
			var anonymous interface{}
			list = append(list, &anonymous)
		}
	}

	return list, nil
}

func getFields(v reflect.Value) (map[string]reflect.Value, error) {
	ve := reflect.Indirect(v)
	vt := indirect(ve.Type())
	fields := make(map[string]reflect.Value)
	for i := 0; i < ve.NumField(); i++ {
		fv := ve.Field(i)
		ft := vt.Field(i)
		tag := getTag(ft)

		fvt := indirect(fv.Type())
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			if !fv.CanInterface() {
				return nil, errNotSettable
			}

			nfv := reflect.New(indirect(fv.Type()))
			fv.Set(nfv)
			fvt = indirect(fvt)
		}

		if fvt.Kind() == reflect.Struct && ft.Anonymous {
			fds, err := getFields(fv)
			if err != nil {
				return nil, err
			}

			for k, v := range fds {
				fields[k] = v
			}
		} else {
			if !fv.CanAddr() || !fv.Addr().CanInterface() {
				return nil, errNotSettable
			}

			fields[tag] = fv.Addr()
		}
	}

	return fields, nil
}

func getTag(f reflect.StructField) string {
	tag, ok := f.Tag.Lookup("db")
	if !ok {
		tag = f.Name
	} else {
		index := strings.Index(tag, ",")
		if index > 0 {
			tag = tag[:index]
		}
	}

	return tag
}
func scanBasicRow(rows *sql.Rows, v interface{}) error {
	value := reflect.ValueOf(v)
	elem := reflect.Indirect(value)
	if !elem.CanSet() {
		return errNotSettable
	}

	return rows.Scan(value.Interface())
}

func indirect(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}

func must(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return errInvalidPointer
	}

	value := reflect.ValueOf(v)
	if value.IsNil() {
		return errInvalidPointer
	}

	if !value.IsValid() || !value.CanInterface() {
		return errInvalidPointer
	}

	return nil
}
