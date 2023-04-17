package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrNotAStructure = errors.New("not a structure")
	ErrIntMin        = errors.New("value less than min")
	ErrIntMax        = errors.New("int value more than max")
	ErrIntNotFound   = errors.New("value not found in dictionary")
	ErrStrLen        = errors.New("length of string not equal to ethalon")
	ErrStrRxp        = errors.New("value not succeeded to regexp")
	ErrStrNotFound   = errors.New("string value not found in dictionary")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	str := ""
	for _, ve := range v {
		str += ve.Field + ": " + ve.Err.Error() + "\r\n"
	}
	return strings.TrimRight(str, "\r\n")
}

func (v ValidationErrors) Is(x error) bool {
	for _, ve := range v {
		if errors.Is(ve.Err, x) {
			return true
		}
	}
	return false
}

func Validate(v interface{}) error {
	errs := make(ValidationErrors, 0)
	validateX(v, "", &errs)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

//nolint:exhaustive
func validateX(u interface{}, nameParent string, errs *ValidationErrors) {
	var chkErrArr []error
	t := reflect.TypeOf(u)
	if t.Kind() != reflect.Struct {
		*errs = append(*errs, ValidationError{Field: nameParent, Err: ErrNotAStructure})
		return
	}
	v := reflect.ValueOf(u)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fnameFull := strings.TrimLeft(nameParent+"."+f.Name, ".")
		tag := f.Tag.Get("validate")
		if tag == "" {
			continue
		}
		tags := strings.Split(tag, "|")
		fv := v.Field(i)
		switch f.Type.Kind() {
		case reflect.Struct:
			if tag == "nested" {
				validateX(fv.Interface(), fnameFull, errs)
			}
		case reflect.String:
			strv := make([]string, 1)
			strv[0] = fv.String()
			chkErrArr = stringValidate(strv, tags)
		case reflect.Int:
			intv := make([]int, 1)
			intv[0] = fv.Interface().(int)
			chkErrArr = intValidate(intv, tags)
		case reflect.Slice:
			switch f.Type.Elem().Kind() {
			case reflect.String:
				chkErrArr = stringValidate(fv.Interface().([]string), tags)
			case reflect.Int:
				chkErrArr = intValidate(fv.Interface().([]int), tags)
			}
		default:
			continue
		}
		if len(chkErrArr) > 0 {
			for _, chkErr := range chkErrArr {
				*errs = append(*errs, ValidationError{Field: fnameFull, Err: chkErr})
			}
		}
	}
}

func intValidate(intArray []int, tags []string) []error {
	var fCheck func(x int, info string) error
	outErr := make([]error, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "min":
			fCheck = intCheckMin
		case "max":
			fCheck = intCheckMax
		case "in":
			fCheck = intCheckDict
		}
		for _, v := range intArray {
			errChk := fCheck(v, tagParsed[1])
			if errChk != nil {
				outErr = append(outErr, errChk)
			}
		}
	}
	return outErr
}

func stringValidate(strArray []string, tags []string) []error {
	var fCheck func(x string, info string) error
	outErr := make([]error, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "len":
			fCheck = stringCheckLen
		case "regexp":
			fCheck = stringCheckRxp
		case "in":
			fCheck = stringCheckDict
		}
		for _, v := range strArray {
			errChk := fCheck(v, tagParsed[1])
			if errChk != nil {
				outErr = append(outErr, errChk)
			}
		}
	}
	return outErr
}
