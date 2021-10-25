package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	AppInvalidLenValue struct {
		Version string `validate:"len:a"`
	}

	UserInvalidRegexp struct {
		Email string `validate:"regexp:)))"`
		Age   int    `validate:"min:a|max:b"`
	}

	UserInvalidAge struct {
		Age int `validate:"min:a|max:b"`
	}

	ResponseInvalidValuesList struct {
		Code int `validate:"in:a,b,c"`
	}

	App struct {
		Version string `validate:"len:5|regexp:^\\d+"`
	}

	Token struct {
		Header  []byte
		Payload []byte `validate:"len:11"`
	}

	NestedUser struct {
		ID       string
		User     *User  `validate:"nested"`
		password string `validate:"min:8|max:30"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"app version length",
			App{
				Version: "12345",
			},
			nil,
		},

		{
			"age in lower range interval",
			User{
				Age:   18,
				Role:  UserRole("admin"),
				Email: "test@otus.ru",
			},
			nil,
		},
		{
			"age in upper range interval",
			User{
				Age:   50,
				Role:  UserRole("stuff"),
				Email: "test@otus.ru",
			},
			nil,
		},
		{
			"email regexp",
			User{
				Age:   20,
				Email: "test@otus.ru",
				Role:  UserRole("stuff"),
			},
			nil,
		},
		{
			"phones slice item length",
			User{
				Age:   20,
				Role:  UserRole("admin"),
				Email: "test@otus.ru",
				Phones: []string{
					"01234567890",
					"09876543210",
				},
			},
			nil,
		},

		{
			"response code 200 in list",
			Response{
				Code: 200,
			},
			nil,
		},
		{
			"response code 404 in list",
			Response{
				Code: 404,
			},
			nil,
		},
		{
			"response code 500 in list",
			Response{
				Code: 500,
			},
			nil,
		},
		{
			"response body not need to validate",
			Response{
				Code: 200,
				Body: "OK",
			},
			nil,
		},

		{
			"nested struct",
			NestedUser{
				User: &User{
					Age:   18,
					Role:  UserRole("admin"),
					Email: "test@otus.ru",
				},
				password: "test",
			},
			nil,
		},

		{
			"slice bytes length",
			Token{
				Payload: []byte("abcdefghij"),
			},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestInvalidValidatorTag(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"len parse error",
			AppInvalidLenValue{
				Version: "123",
			},
			ErrParseValidator,
		},
		{
			"values list parse error",
			ResponseInvalidValuesList{
				Code: 500,
			},
			ErrParseValidator,
		},
		{
			"regexp parse error",
			UserInvalidRegexp{
				Email: "otus@test.ru",
			},
			ErrInvalidRegexp,
		},
		{
			"range parse error",
			UserInvalidAge{
				Age: -25,
			},
			ErrParseValidator,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"version length less",
			App{Version: "1234"},
			ValidationErrors{
				{Field: "Version", Err: ErrStringLen},
			},
		},
		{
			"version regexp not matched",
			App{Version: "a1234"},
			ValidationErrors{
				{Field: "Version", Err: ErrRegexpNotMatched},
			},
		},
		{
			"version length greater",
			App{Version: "123456"},
			ValidationErrors{
				{Field: "Version", Err: ErrStringLen},
			},
		},

		{
			"age value less",
			User{Age: 17},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
		{
			"age value greater",
			User{Age: 51},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMax},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
		{
			"email regexp not matched",
			User{Email: "test"},
			ValidationErrors{
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
		{
			"phones slice empty",
			User{},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
		{
			"phones slice item length less",
			User{Phones: []string{"0123456789"}},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
				{Field: "Phones", Err: ErrStringLen},
			},
		},
		{
			"role value not in list",
			User{Role: UserRole("test")},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},

		{
			"response code not in list",
			Response{Code: 301},
			ValidationErrors{
				{Field: "Code", Err: ErrValueNotInList},
			},
		},

		{
			"nested struct",
			NestedUser{
				User: &User{
					Age:   17,
					Role:  UserRole("user"),
					Email: "test",
				},
			},
			ValidationErrors{
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrRegexpNotMatched},
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
