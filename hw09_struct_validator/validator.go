package hw09structvalidator

import (
	"errors"
	"fmt"
	"log/slog"
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
	ErrInvalidIncomingValue = errors.New("invalid incoming value, expected structure")
	ErrInvalidValidationTag = errors.New("invalid validation tag")
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
	var err string

	for i := 0; i < len(v); i++ {
		err += fmt.Sprintf("%s errors: %s;", v[i].Field, v[i].Err.Error())
	}

	return err
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
					err = fmt.Errorf("validate %s value err: %w", fieldName, err)
					slog.Error(err.Error())
					return err
				}

				if fieldErrors != nil {
					validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fieldErrors})
				}

				continue
			}

			if kind == reflect.Array {
				for j := 0; j < reflectVal.Len(); j++ {
					fieldErrors, err := validate(reflectVal.Index(j).String(), validators)
					if err != nil {
						err = fmt.Errorf("validate %s %d elem err: %w", fieldName, j, err)
						slog.Error(err.Error())
						return err
					}

					if fieldErrors != nil {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fieldErrors})
					}
				}

				continue
			}
		}

	}

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
			err := validateLen(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}

			err = validateRegexp(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}

			err = validateIn(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}
		case int:
			err := validateMin(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}

			err = validateMax(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}

			err = validateIn(validator, val, &fieldErrors)
			if err != nil {
				return nil, err
			}
		}
	}

	return fieldErrors, nil
}

func validateLen(validator []string, value string, fieldErrors *error) error {
	if validator[0] == "len" {
		length, err := strconv.Atoi(validator[1])
		if err != nil {
			return fmt.Errorf("parse len to int err: %w", err)
		}

		if utf8.RuneCountInString(value) != length {
			if *fieldErrors != nil {
				*fieldErrors = fmt.Errorf("%w, %w", *fieldErrors, ErrIncorrectStrLength)
				return nil
			}

			*fieldErrors = ErrIncorrectStrLength
		}
	}

	return nil
}

func validateRegexp(validator []string, value string, fieldErrors *error) error {
	if validator[0] == "regexp" {
		re, err := regexp.Compile(validator[1])
		if err != nil {
			return fmt.Errorf("regexp.Compile err: %w", err)
		}

		ok := re.MatchString(value)
		if !ok {
			if *fieldErrors != nil {
				*fieldErrors = fmt.Errorf("%w, %w", *fieldErrors, ErrIncorrectStrContent)
				return nil
			}

			*fieldErrors = ErrIncorrectStrContent
		}
	}

	return nil
}

func validateMin(validator []string, value int, fieldErrors *error) error {
	if validator[0] == "min" {
		min, err := strconv.Atoi(validator[1])
		if err != nil {
			return fmt.Errorf("parse min to int err: %w", err)
		}

		if value < min {
			if *fieldErrors != nil {
				*fieldErrors = fmt.Errorf("%w, %w", *fieldErrors, ErrMinNotMet)
				return nil
			}

			*fieldErrors = ErrMinNotMet
		}
	}

	return nil
}

func validateMax(validator []string, value int, fieldErrors *error) error {
	if validator[0] == "max" {
		max, err := strconv.Atoi(validator[1])
		if err != nil {
			return fmt.Errorf("parse max to int err: %w", err)
		}

		if value > max {
			if *fieldErrors != nil {
				*fieldErrors = fmt.Errorf("%w, %w", *fieldErrors, ErrMaxNotMet)
				return nil
			}

			*fieldErrors = ErrMaxNotMet
		}
	}

	return nil
}

func validateIn[T validationType](validator []string, value T, fieldErrors *error) error {
	if validator[0] == "in" {
		expectedValues := strings.Split(validator[1], ",")
		var isSuccess bool

		switch val := any(value).(type) {
		case string:
			for i := 0; i < len(expectedValues); i++ {
				if expectedValues[i] == val {
					isSuccess = true
					break
				}
			}
		case int:
			for i := 0; i < len(expectedValues); i++ {
				expectedValue, err := strconv.Atoi(expectedValues[i])
				if err != nil {
					return fmt.Errorf("parse expected value to int err: %w", err)
				}

				if expectedValue == val {
					isSuccess = true
					break
				}
			}
		default:
			return nil
		}

		if !isSuccess {
			if *fieldErrors != nil {
				*fieldErrors = fmt.Errorf("%w, %w", *fieldErrors, ErrUnexpectedValue)
				return nil
			}

			*fieldErrors = ErrUnexpectedValue
		}
	}

	return nil
}
