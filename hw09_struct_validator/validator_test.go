package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		Usr     User   `validate:"nested"`
	}

	App2 struct {
		Version string `validate:"len:5"`
		Usr     User
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			"---",
			ErrNotAStructure,
		},
		{
			User{ID: "-", Age: 22, Email: "a@b.c", Role: "stuff"},
			ErrStrLen,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 22, Email: "a@b.c",
				Role: "stuff", Phones: []string{"11111111111", "2"},
			},
			ErrStrLen,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 22, Email: "a@b@c",
				Role: "stuff",
			},
			ErrStrRxp,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 22, Email: "a@b.c",
				Role: "-",
			},
			ErrStrNotFound,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 15, Email: "a@b.c",
				Role: "admin",
			},
			ErrIntMin,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 51, Email: "a@b.c",
				Role: "stuff",
			},
			ErrIntMax,
		},
		{
			User{
				ID: "123456789012345678901234567890666666", Age: 50, Email: "a@b.c",
				Role: "stuff",
			},
			nil,
		},
		{
			Response{Code: 666, Body: "stuff"},
			ErrIntNotFound,
		},
		{
			Token{
				Header: []byte{0, 34, 234, 234}, Payload: []byte{0, 34, 234, 234},
				Signature: []byte{0, 34, 234, 234},
			},
			nil,
		},
		{
			App{Version: "XXXXX", Usr: User{
				ID: "123456789012345678901234567890666666", Age: 51, Email: "a@b.c",
				Role: "stuff",
			}},
			ErrIntMax,
		},
		{
			App2{Version: "XXXXX", Usr: User{
				ID: "123456789012345678901234567890666666", Age: 51, Email: "a@b.c",
				Role: "stuff",
			}},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Truef(t, errors.Is(err, tt.expectedErr), "got err [%v] but expected err [%v]", err, tt.expectedErr)
			_ = tt
		})
	}
}
