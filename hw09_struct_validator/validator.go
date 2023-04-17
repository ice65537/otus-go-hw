package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

func Validate(v interface{}) error {
	errs := make(ValidationErrors, 0)
	validateX(v, "", &errs)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateX(u interface{}, nameParent string, errs *ValidationErrors) {
	var chkErrArr []error
	fmt.Println("Structure->" + nameParent)
	t := reflect.TypeOf(u)
	if t.Kind() != reflect.Struct {
		*errs = append(*errs, ValidationError{Field: nameParent, Err: ErrNotAStructure})
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
		fmt.Println("Field--->" + f.Name)
		switch f.Type.Kind() {
		case reflect.Struct:
			validateX(fv.Interface(), strings.TrimLeft(nameParent+"."+f.Name, "."), errs)
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
				*errs = append(*errs, ValidationError{Field: strings.TrimLeft(nameParent+"."+f.Name, "."), Err: chkErr})
			}
		}
	}
}

func intValidate(intArray []int, tags []string) []error {
	var fCheck func(x int) error
	outErr := make([]error, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "min":
			fCheck = func(x int) error {
				min, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err
				}
				if x < min {
					fmt.Printf("Value [%d] less than min=[%d]\r\n", x, min)
					return ErrIntMin
				}
				return nil
			}
		case "max":
			fCheck = func(x int) error {
				max, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err
				}
				if x > max {
					fmt.Printf("Value [%d] more than max=[%d]\r\n", x, max)
					return ErrIntMax
				}
				return nil
			}
		case "in":
			fCheck = func(x int) error {
				var intSet map[int]struct{}
				intSetStr := strings.Split(tagParsed[1], ",")
				intSet = make(map[int]struct{})
				for _, v := range intSetStr {
					idx, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
					intSet[idx] = struct{}{}
				}
				_, ok := intSet[x]
				if !ok {
					fmt.Printf("Value [%d] not found in dictionary[%s]\r\n", x, tagParsed[1])
					return ErrIntNotFound
				}
				return nil
			}
		}
		for _, v := range intArray {
			errChk := fCheck(v)
			if errChk != nil {
				outErr = append(outErr, errChk)
			}
		}
	}
	return outErr
}

func stringValidate(strArray []string, tags []string) []error {
	var fCheck func(x string) error
	outErr := make([]error, 0)
	for _, tag := range tags {
		tagParsed := strings.Split(tag, ":")
		switch tagParsed[0] {
		case "len":
			fCheck = func(x string) error {
				exactLength, err := strconv.Atoi(tagParsed[1])
				if err != nil {
					return err
				}
				if len(x) != exactLength {
					fmt.Printf("Length of string [%s] not equal to [%d]\r\n", x, exactLength)
					return ErrStrLen
				}
				return nil
			}
		case "regexp":
			fCheck = func(x string) error {
				rex, err := regexp.Compile(tagParsed[1])
				if err != nil {
					return err
				}
				if !rex.MatchString(x) {
					fmt.Printf("Value [%s] not succeeded to regexp [%s]\r\n", x, tagParsed[1])
					return ErrStrRxp
				}
				return nil
			}
		case "in":
			fCheck = func(x string) error {
				var strSet map[string]struct{}
				strSetSlice := strings.Split(tagParsed[1], ",")
				strSet = make(map[string]struct{})
				for _, idx := range strSetSlice {
					strSet[idx] = struct{}{}
				}
				_, ok := strSet[x]
				if !ok {
					fmt.Printf("Value [%s] not found in dictionary[%s]\r\n", x, tagParsed[1])
					return ErrStrNotFound
				}
				return nil
			}
		}
		for _, v := range strArray {
			errChk := fCheck(v)
			if errChk != nil {
				outErr = append(outErr, errChk)
			}
		}
	}
	return outErr
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
