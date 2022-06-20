package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	InvalidValidate struct {
		Name string `validate:"len:2:4"`
	}

	InvalidValidatetReg struct {
		Name string `validate:"regexp:\\\\\\\\\\"`
	}

	InvalidValidateInt struct {
		Name string `validate:"len:sdf"`
	}

	PrivateValidateFields struct {
		name string `validate:"len:5"`
		age  int    `validate:"min:18|max:50"`
	}

	NotSupportedValidateFields struct {
		Name struct{} `validate:"len:5"`
	}

	DoubleValidate struct {
		Name string `validate:"len:2|len:3"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{App{Version: "1234"}, ValidationErrors{ValidationError{"Version", ErrLen}}},
		{App{Version: "12345"}, nil},
		{App{Version: "123456"}, ValidationErrors{ValidationError{Field: "Version", Err: ErrLen}}},

		{Response{Code: 199, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 200, Body: ""}, nil},
		{Response{Code: 201, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{Response{Code: 403, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 404, Body: ""}, nil},
		{Response{Code: 405, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{Response{Code: 499, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		{Response{Code: 500, Body: ""}, nil},
		{Response{Code: 501, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{User{
			ID:     "",
			Name:   "",
			Age:    0,
			Email:  "",
			Role:   "",
			Phones: []string{""},
			meta:   json.RawMessage(""),
		}, ValidationErrors([]ValidationError{
			{Field: "ID", Err: ErrLen},
			{Field: "Age", Err: ErrMin},
			{Field: "Email", Err: ErrReg},
			{Field: "Role", Err: ErrIn},
			{Field: "Phones", Err: ErrLen},
		})},
		{User{
			ID:     "123456789123456789123456789123456789",
			Name:   "",
			Age:    19,
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"12345678912"},
		}, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &ValidationErrors{})
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestInvalidValidate(t *testing.T) {
	err := Validate(InvalidValidate{Name: ""})
	require.ErrorAs(t, err, &ValidatorError{})
}

func TestInvalidValidatetReg(t *testing.T) {
	err := Validate(InvalidValidatetReg{Name: ""})
	require.ErrorAs(t, err, &ValidatorError{})
}

func TestInvalidValidateInt(t *testing.T) {
	err := Validate(InvalidValidateInt{Name: ""})
	require.ErrorAs(t, err, &ValidatorError{})
}

func TestDoubleValidate(t *testing.T) {
	err := Validate(DoubleValidate{Name: ""})
	require.ErrorAs(t, err, &ValidatorError{})
}

func TestPrivateValidateFields(t *testing.T) {
	err := Validate(PrivateValidateFields{name: "", age: 30})
	require.NoError(t, err)
}

func TestNotAStruct(t *testing.T) {
	err := Validate(2)
	require.ErrorAs(t, err, &ValidatorError{})
}

func TestNotSupportedValidateFields(t *testing.T) {
	err := Validate(NotSupportedValidateFields{Name: struct{}{}})
	require.NoError(t, err)
	err = Validate(Token{})
	require.NoError(t, err)
}
