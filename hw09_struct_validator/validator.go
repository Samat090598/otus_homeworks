package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Field string
	Err   error
}

// System errors
var (
	ErrInvalidIncomingValue  = errors.New("invalid incoming value, expected structure")
	ErrInvalidValidationTag  = errors.New("invalid validation tag")
	ErrIncorrectLengthFormat = errors.New("incorrect length format")

	ErrIncorrectMinFormat = errors.New("incorrect min format")
	ErrIncorrectMaxFormat = errors.New("incorrect max format")
)

// Validation errors
var (
	ErrIncorrectStrLength  = errors.New("the length of the string does not match the required one")
	ErrIncorrectStrContent = errors.New("the string does not meet the requirements of the regular expression")
	ErrUnexpectedValue     = errors.New("the value does not match the expected values")

	ErrMinNotMet = errors.New("the value does not correspond to the min")
	ErrMaxNotMet = errors.New("the value does not correspond to the max")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

type validationType interface {
	string | int
}

func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return ErrInvalidIncomingValue
	}

	st := reflect.TypeOf(v)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)

		if validateTag, ok := field.Tag.Lookup("validate"); ok {
			var (
				tp         = field.Type
				kind       = tp.Kind()
				reflectVal = reflect.ValueOf(field)
				fieldName  = field.Name
				validators = strings.Split(validateTag, "|")
			)

			if kind == reflect.String || kind == reflect.Int {
				fieldErrors, err := validate(reflectVal.String(), validators)
				if err != nil {
					return fmt.Errorf("validate %s value err: %w", fieldName, err)
				}

				if fieldErrors != nil {
					validationErrors = append(validationErrors, ValidationError{Field: fieldName})
				}

				continue
			}

			if kind == reflect.Array {
				for j := 0; j < reflectVal.Len(); j++ {
					fieldErrors, err := validate(reflectVal.Index(j).String(), validators)
					if err != nil {
						return fmt.Errorf("validate %s %d elem err: %w", fieldName, j, err)
					}

					if fieldErrors != nil {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName})
					}
				}

				continue
			}
		}

	}

	// Place your code here.
	return validationErrors
}

func validate[T validationType](value T, validators []string) (error, error) {
	var fieldErrors error

	for i := 0; i < len(validators); i++ {
		validator := strings.SplitN(validators[i], ":", 2)
		if len(validator) < 2 {
			return nil, ErrInvalidValidationTag
		}

		switch val := any(value).(type) {
		case string:
			if validator[0] == "len" {
				length, err := strconv.Atoi(validator[1])
				if err != nil {
					return nil, ErrIncorrectLengthFormat
				}

				if utf8.RuneCountInString(val) != length {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrIncorrectStrLength)
				}

				continue
			}

			if validator[0] == "regexp" {
				re, err := regexp.Compile(validator[1])
				if err != nil {
					return nil, fmt.Errorf("regexp.Compile err: %w", err)
				}

				ok := re.MatchString(validator[1])
				if !ok {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrIncorrectStrContent)
				}

				continue
			}

			if validator[0] == "in" {
				expectedValues := strings.Split(validator[1], ",")
				var isSuccess bool
				for j := 0; j < len(expectedValues); j++ {
					if expectedValues[i] == val {
						isSuccess = true
						break
					}
				}

				if !isSuccess {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrUnexpectedValue)
				}

				continue
			}
		case int:
			if validator[0] == "min" {
				min, err := strconv.Atoi(validator[1])
				if err != nil {
					return nil, ErrIncorrectMinFormat
				}

				if val < min {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrMinNotMet)
					continue
				}
			}

			if validator[0] == "max" {
				min, err := strconv.Atoi(validator[1])
				if err != nil {
					return nil, ErrIncorrectMaxFormat
				}

				if val < min {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrMaxNotMet)
					continue
				}
			}

			if validator[0] == "in" {
				expectedValues := strings.Split(validator[1], ",")
				var isSuccess bool
				for j := 0; j < len(expectedValues); j++ {
					expectedValue, err := strconv.Atoi(expectedValues[i])
					if err != nil {
						return nil, fmt.Errorf("parse expected value to int err: %w", err)
					}

					if expectedValue == val {
						isSuccess = true
						break
					}
				}

				if !isSuccess {
					fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrUnexpectedValue)
				}

				continue
			}
		}
	}

	return fieldErrors, nil
}