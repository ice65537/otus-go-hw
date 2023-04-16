package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

func Validate(v interface{}) error {
	errs := make(ValidationErrors, 0)
	validateX(v, "", &errs)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateX(u interface{}, nameParent string, errs *ValidationErrors) {
	var chkMsgArr []string
	t := reflect.TypeOf(u)
	if t.Kind() != reflect.Struct {
		*errs = append(*errs, ValidationError{Field: nameParent, Err: fmt.Errorf("%s is not a structure", t.Kind())})
		return
	}
	v := reflect.ValueOf(u)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("validate")
		if tag == "" {
			continue
		}
		tags := strings.Split(tag, "|")
		fv := v.Field(i)
		fmt.Println(f.Name)
		switch f.Type.Kind() {
		case reflect.Struct:
			validateX(fv.Interface(), strings.TrimLeft(nameParent+"."+f.Name, "."), errs)
		case reflect.String:
			strv := make([]string, 1)
			strv[0] = fv.String()
			chkMsgArr = stringValidate(strv, tags)
		case reflect.Int:
			intv := make([]int, 1)
			intv[0] = fv.Interface().(int)
			chkMsgArr = intValidate(intv, tags)
		case reflect.Slice:
			switch f.Type.Elem().Kind() {
			case reflect.String:
				chkMsgArr = stringValidate(fv.Interface().([]string), tags)
			case reflect.Int:
				chkMsgArr = intValidate(fv.Interface().([]int), tags)
			}
		default:
			continue
		}
		if len(chkMsgArr) > 0 {
			for _, chkMsg := range chkMsgArr {
				*errs = append(*errs, ValidationError{Field: strings.TrimLeft(nameParent+"."+f.Name, "."), Err: errors.New(chkMsg)})
			}
		}
	}
}

func intValidate(intArray []int, tags []string) []string {
	var fCheck func(x int) string
	outStr := make([]string, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "min":
			fCheck = func(x int) string {
				min, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err.Error()
				}
				if x < min {
					return fmt.Sprintf("Value [%d] less than min=[%d]", x, min)
				}
				return ""
			}
		case "max":
			fCheck = func(x int) string {
				max, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err.Error()
				}
				if x > max {
					return fmt.Sprintf("Value [%d] more than max=[%d]", x, max)
				}
				return ""
			}
		case "in":
			fCheck = func(x int) string {
				var intSet map[int]struct{}
				intSetStr := strings.Split(tagParsed[1], ",")
				intSet = make(map[int]struct{})
				for _, v := range intSetStr {
					idx, err := strconv.Atoi(v)
					if err != nil {
						return err.Error()
					}
					intSet[idx] = struct{}{}
				}
				_, ok := intSet[x]
				if !ok {
					return fmt.Sprintf("Value [%d] not found in dictionary[%s]", x, tagParsed[1])
				}
				return ""
			}
		}
		for _, v := range intArray {
			errStr := fCheck(v)
			if errStr != "" {
				outStr = append(outStr, errStr)
			}
		}
	}
	return outStr
}

func stringValidate(strArray []string, tags []string) []string {
	var fCheck func(x string) string
	outStr := make([]string, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "len":
			fCheck = func(x string) string {
				exactLength, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err.Error()
				}
				if len(x) != exactLength {
					return fmt.Sprintf("Length of string [%s] not equal to [%d]", x, exactLength)
				}
				return ""
			}
		case "regexp":
			fCheck = func(x string) string {
				rex, err := regexp.Compile(tagParsed[1])
				if err != nil {
					return err.Error()
				}
				if !rex.MatchString(x) {
					return fmt.Sprintf("Value [%s] not succeeded to regexp [%s]", x, tagParsed[1])
				}
				return ""
			}
		case "in":
			fCheck = func(x string) string {
				var strSet map[string]struct{}
				strSetSlice := strings.Split(tagParsed[1], ",")
				strSet = make(map[string]struct{})
				for _, idx := range strSetSlice {
					strSet[idx] = struct{}{}
				}
				_, ok := strSet[x]
				if !ok {
					return fmt.Sprintf("Value [%s] not found in dictionary[%s]", x, tagParsed[1])
				}
				return ""
			}
		}
		for _, v := range strArray {
			errStr := fCheck(v)
			if errStr != "" {
				outStr = append(outStr, errStr)
			}
		}
	}
	return outStr
}

/*Необходимо реализовать следующие валидаторы:
- Для строк:
    * `len:32` - длина строки должна быть ровно 32 символа;
    * `regexp:\\d+` - согласно регулярному выражению строка должна состоять из цифр
    (`\\` - экранирование слэша);
    * `in:foo,bar` - строка должна входить в множество строк {"foo", "bar"}.
- Для чисел:
    * `min:10` - число не может быть меньше 10;
    * `max:20` - число не может быть больше 20;
    * `in:256,1024` - число должно входить в множество чисел {256, 1024};
- Для слайсов валидируется каждый элемент слайса.

_При желании можно дополнительно добавить парочку новых правил (на ваше усмотрение)._

Допускается комбинация валидаторов по логическому "И" с помощью `|`, например:
* `min:0|max:10` - число должно находится в пределах [0, 10];
* `regexp:\\d+|len:20` - строка должна состоять из цифр и иметь длину 20.*/
