package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Animal struct {
		Name string `validate:"test"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "1",
				Name:   "John",
				Age:    17,
				Email:  "john.dow@gmail.com",
				Role:   "guest",
				Phones: []string{"11111111111", "2222222222", "333333333"},
			},
			expectedErr: fmt.Errorf("{ID %w}; {Age %w}; {Email %w}; {Role %w}; {Phones 1 elem errors: %w, 2 elem errors: %w}; ",
				ErrIncorrectStrLength, ErrMinNotMet, ErrIncorrectStrContent, ErrUnexpectedValue, ErrIncorrectStrLength, ErrIncorrectStrLength),
		},
		{
			in:          App{Version: "123456"},
			expectedErr: fmt.Errorf("{Version %w}; ", ErrIncorrectStrLength),
		},
		{
			in: Token{
				Header:    make([]byte, 0),
				Payload:   make([]byte, 0),
				Signature: make([]byte, 0),
			},
			expectedErr: errors.New(""),
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: errors.New(""),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorAs(t, err, &ValidationErrors{})
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}

	t.Run("invalid incoming value", func(t *testing.T) {
		err := Validate("test")
		require.EqualError(t, err, ErrInvalidIncomingValue.Error())
	})

	t.Run("invalid validation tag", func(t *testing.T) {
		err := Validate(Animal{Name: "Dog"})
		require.ErrorIs(t, err, ErrInvalidValidationTag)
	})
}
