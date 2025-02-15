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
)

// Validation errors
var (
	ErrIncorrectStrLength  = errors.New("the length of the string does not match the required one")
	ErrIncorrectStrContent = errors.New("the string does not meet the requirements of the regular expression")
	ErrIncorrectStr        = errors.New("the string does not match the expected values")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return ErrInvalidIncomingValue
	}

	st := reflect.TypeOf(v)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)

		if validate, ok := field.Tag.Lookup("validate"); ok {
			var (
				tp         = field.Type
				kind       = tp.Kind()
				reflectVal = reflect.ValueOf(field)
				fieldName  = field.Name
				validators = strings.Split(validate, "|")
			)

			if kind == reflect.String {
				fieldErrors, err := validateString(reflectVal.String(), validators)
				if err != nil {
					return fmt.Errorf("validate %s value err: %w", fieldName, err)
				}

				if fieldErrors != nil {
					validationErrors = append(validationErrors, ValidationError{Field: fieldName})
				}

				continue
			}

			if kind == reflect.Int {
				//...
				continue
			}

			if kind == reflect.Array {
				if tp.Elem().Kind() == reflect.String {
					for j := 0; j < reflectVal.Len(); j++ {
						fieldErrors, err := validateString(reflectVal.Index(j).String(), validators)
						if err != nil {
							return fmt.Errorf("validate %s %d elem err: %w", fieldName, j, err)
						}

						if fieldErrors != nil {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName})
						}
					}

					continue
				}

				if tp.Elem().Kind() == reflect.Int {

				}
			}
		}

	}

	// Place your code here.
	return validationErrors
}

func validateString(value string, validators []string) (error, error) {
	var fieldErrors error

	for i := 0; i < len(validators); i++ {
		validator := strings.SplitN(validators[i], ":", 2)
		if len(validator) < 2 {
			return nil, ErrInvalidValidationTag
		}

		if validator[0] == "len" {
			length, err := strconv.Atoi(validator[1])
			if err != nil {
				return nil, ErrIncorrectLengthFormat
			}

			if utf8.RuneCountInString(value) != length {
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
				if expectedValues[i] == value {
					isSuccess = true
					break
				}
			}

			if !isSuccess {
				fieldErrors = fmt.Errorf("%w, %w", fieldErrors, ErrIncorrectStr)
			}

			continue
		}
	}

	return fieldErrors, nil
}
