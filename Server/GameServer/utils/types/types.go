package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------
// Conversion interfaces

// TypeMarshaller is implemented by any value that has a Marshal method
// This converter is used to convert the value to it string representation
type TypeMarshaller interface {
	Marshal() (string, error)
}

// TypeUnmarshaller is implemented by any value that has an Unmarshal method
// This converter is used to convert a string to your value representation of that string
type TypeUnmarshaller interface {
	Unmarshal(string) error
}

// NoUnmarshalFuncError is the custom error type to be raised in case there is no unmarshal function defined on type
type NoUnmarshalFuncError struct {
	msg string
}

func (e NoUnmarshalFuncError) Error() string {
	return e.msg
}

// NoMarshalFuncError is the custom error type to be raised in case there is no marshal function defined on type
type NoMarshalFuncError struct {
	ty reflect.Type
}

func (e NoMarshalFuncError) Error() string {
	return "No known conversion from " + e.ty.String() + " to string, " + e.ty.String() + " does not implement TypeMarshaller nor Stringer"
}

// --------------------------------------------------------------------------
// Conversion helpers

func ToString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", inValue.Uint()), nil
	case reflect.Float32:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}

func ToBool(in interface{}) (bool, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := inValue.String()
		s = strings.TrimSpace(s)
		if strings.EqualFold(s, "yes") {
			return true, nil
		} else if strings.EqualFold(s, "no") || s == "" {
			return false, nil
		} else {
			return strconv.ParseBool(s)
		}
	case reflect.Bool:
		return inValue.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := inValue.Int()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := inValue.Uint()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		f := inValue.Float()
		if f != 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to bool")
}

func ToInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		out := strings.SplitN(s, ".", 2)
		return strconv.ParseInt(out[0], 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func ToUint(in interface{}) (uint64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}

		// support the float input
		if strings.Contains(s, ".") {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, err
			}
			return uint64(f), nil
		}
		return strconv.ParseUint(s, 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inValue.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to uint")
}

func ToFloat(in interface{}) (float64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func SetField(field reflect.Value, value string, omitEmpty bool) error {
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Interface().(type) {
	case string:
		s, err := ToString(value)
		if err != nil {
			return err
		}
		field.SetString(s)
	case bool:
		b, err := ToBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case int, int8, int16, int32, int64:
		i, err := ToInt(value)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case uint, uint8, uint16, uint32, uint64:
		ui, err := ToUint(value)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case float32, float64:
		f, err := ToFloat(value)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("type not in [string, boolint, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64], got type %s", reflect.TypeOf(field.Interface()).String())
	}
	return nil
}
