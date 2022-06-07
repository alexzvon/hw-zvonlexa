package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var structTypeOf = map[reflect.Kind]struct{}{
	reflect.String: {},
	reflect.Int:    {},
	reflect.Slice:  {},
}

var (
	ErrLen = errors.New("неверная длинна")
	ErrMin = errors.New("число меньше минимального значения")
	ErrMax = errors.New("число больше максимального значения")
	ErrIn  = errors.New("не верное значение")
	ErrReg = errors.New("значение не соответствует регулярному выражению")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var strBuild strings.Builder

	for _, validError := range ve {
		strBuild.WriteString(concat("Поле \"", validError.Field, "\": ", validError.Err.Error(), "\n"))
	}

	return strBuild.String()
}

type ValidatorError struct {
	Msg string
	Err error
}

func (ve ValidatorError) Error() string {
	return ve.Msg
}

type validator struct {
	errors      ValidationErrors
	duplicate   map[string]struct{}
	structField reflect.StructField
}

func (v *validator) validatorSplitTag(str string) (string, string, error) {
	sTag := strings.Split(str, ":")

	if len(sTag) != 2 {
		return "", "", ValidatorError{
			Msg: fmt.Sprintf(
				"Поле %q - не верный формат, ожидалось [name:value], но получено %v",
				v.structField.Name,
				sTag,
			),
		}
	}

	keyTag := sTag[0]
	if _, ok := v.duplicate[keyTag]; ok {
		return "", "", ValidatorError{
			Msg: fmt.Sprintf(
				"В поле %q - уже имеется ограничение %q",
				v.structField.Name,
				keyTag,
			),
		}
	}

	return keyTag, sTag[1], nil
}

func (v *validator) validatorString(valueField string, vList []string) error {
	for _, list := range vList {
		keyTag, valueTag, errT := v.validatorSplitTag(list)
		if errT != nil {
			return errT
		}

		v.duplicate[keyTag] = struct{}{}

		switch keyTag {
		case "len":
			if err := v.strLen(valueField, keyTag, valueTag); err != nil {
				return err
			}
		case "in":
			if err := v.strIn(valueField, keyTag, valueTag); err != nil {
				return err
			}
		case "regexp":
			if err := v.strRegexp(valueField, keyTag, valueTag); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *validator) strLen(valueField, keyTag, valueTag string) error {
	num, err := strconv.Atoi(valueTag)
	if err != nil {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"В поле %q ожидалось числовое значение для ограничения %q, но полученно значение %v",
				v.structField.Name,
				keyTag,
				valueTag,
			),
			Err: err,
		}
	}

	if len(valueField) != num {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrLen),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func (v *validator) strIn(valueField, keyTag, valueTag string) error {
	var in bool

	values := strings.Split(valueTag, ",")

	if len(values) == 0 {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"В ограничении %q нет значения %v, для поиска в поле %q",
				keyTag,
				valueTag,
				v.structField.Name,
			),
		}
	}

	for _, val := range values {
		if val == valueField {
			in = true
		}
	}

	if !in {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrIn),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func (v *validator) strRegexp(valueField, keyTag, valueTag string) error {
	reg, err := regexp.Compile(valueTag)
	if err != nil {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"Поле %q - не верное регулярное выражение для ограничения %q",
				v.structField.Name,
				keyTag,
			),
			Err: err,
		}
	}

	if !reg.MatchString(valueField) {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrReg),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func (v *validator) validatorInt(valueField int64, vList []string) error {
	for _, list := range vList {
		keyTag, valueTag, errT := v.validatorSplitTag(list)
		if errT != nil {
			return errT
		}

		v.duplicate[keyTag] = struct{}{}

		switch keyTag {
		case "min":
			if err := v.intMin(valueField, keyTag, valueTag); err != nil {
				return err
			}
		case "max":
			if err := v.intMax(valueField, keyTag, valueTag); err != nil {
				return err
			}
		case "in":
			if err := v.intIn(valueField, keyTag, valueTag); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *validator) intMin(valueField int64, keyTag, valueTag string) error {
	num, err := strconv.Atoi(valueTag)
	if err != nil {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"В поле %q ожидалось числовое значение для ограничения %q, но полученно значение %v",
				v.structField.Name,
				keyTag,
				valueTag,
			),
			Err: err,
		}
	}

	if valueField < int64(num) {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrMin),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func (v *validator) intMax(valueField int64, keyTag, valueTag string) error {
	num, err := strconv.Atoi(valueTag)
	if err != nil {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"В поле %q ожидалось числовое значение для ограничения %q, но полученно значение %v",
				v.structField.Name,
				keyTag,
				valueTag,
			),
			Err: err,
		}
	}

	if valueField > int64(num) {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrMax),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func (v *validator) intIn(valueField int64, keyTag, valueTag string) error {
	var in bool

	values := strings.Split(valueTag, ",")

	if len(values) == 0 {
		return ValidatorError{
			Msg: fmt.Sprintf(
				"В ограничении %q нет значения %v, для поиска в поле %q",
				keyTag,
				valueTag,
				v.structField.Name,
			),
		}
	}

	for _, item := range values {
		num, err := strconv.Atoi(item)
		if err != nil {
			return ValidatorError{
				Msg: fmt.Sprintf(
					"В поле %q ожидалось числовое значение для ограничения %q, но полученно значение %v",
					v.structField.Name,
					keyTag,
					valueTag,
				),
				Err: err,
			}
		}

		if valueField == int64(num) {
			in = true
		}
	}

	if !in {
		vError := ValidationError{
			Field: v.structField.Name,
			Err:   fmt.Errorf("%w", ErrIn),
		}

		v.errors = append(v.errors, vError)
	}

	return nil
}

func Validate(v interface{}) error {
	valid := validator{errors: make(ValidationErrors, 0)}

	rV := reflect.ValueOf(v)
	if rV.Kind() != reflect.Struct {
		return ValidatorError{
			Msg: fmt.Sprintf("%v - не является структурой", v),
		}
	}

	rT := rV.Type()

	for i := 0; i < rT.NumField(); i++ {
		fV := rV.Field(i)
		field := rT.Field(i)

		validate, ok := tagLookup(fV, field)
		if !ok {
			continue
		}

		valid.duplicate = make(map[string]struct{})
		valid.structField = field
		vList := strings.Split(validate, "|")

		if err := validateKind(fV, vList, &valid); err != nil {
			return err
		}
	}

	if len(valid.errors) != 0 {
		return valid.errors
	}

	return nil
}

func validateKind(fV reflect.Value, vList []string, valid *validator) error {
	switch fV.Kind() { //nolint:exhaustive
	case reflect.String:
		if err := valid.validatorString(fV.String(), vList); err != nil {
			return err
		}
	case reflect.Int:
		if err := valid.validatorInt(fV.Int(), vList); err != nil {
			return err
		}
	case reflect.Slice:
		switch fV.Interface().(type) {
		case []string:
			sV := fV.Interface().([]string)
			for _, val := range sV {
				if err := valid.validatorString(val, vList); err != nil {
					return err
				}
			}
		case []int:
			sV := fV.Interface().([]int)
			for _, val := range sV {
				if err := valid.validatorInt(int64(val), vList); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func tagLookup(refValue reflect.Value, refSF reflect.StructField) (string, bool) {
	if !refValue.CanInterface() {
		return "", false
	}

	if _, ok := structTypeOf[refValue.Kind()]; !ok {
		return "", false
	}

	tag := refSF.Tag

	v, ok := tag.Lookup("validate")
	if !ok {
		return "", false
	}

	return v, true
}

func concat(s ...string) string {
	var builder strings.Builder
	var long int

	for _, v := range s {
		long += len(v)
	}

	builder.Grow(long)

	for _, v := range s {
		builder.WriteString(v)
	}

	return builder.String()
}
