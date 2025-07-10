package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateLogin(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr error
	}{
		{
			name:    "valid",
			login:   "User1999",
			wantErr: nil,
		},
		{
			name:    "invalid: too short",
			login:   "User19",
			wantErr: errInvalidLogin,
		},
		{
			name:    "invalid: empty",
			login:   "",
			wantErr: errInvalidLogin,
		},
		{
			name:    "invalid characters",
			login:   "32#5User",
			wantErr: errInvalidLogin,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogin(tt.login)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid",
			password: "Pass1999@w",
			wantErr:  nil,
		},
		{
			name:     "invalid: too short",
			password: "Pa19@w",
			wantErr:  errInvalidPassTooShort,
		},
		{
			name:     "invalid: no uppercase letters",
			password: "password19@",
			wantErr:  errInvalidPassCase,
		},
		{
			name:     "invalid: no digits",
			password: "passWORD$",
			wantErr:  errInvalidPassNoDigit,
		},
		{
			name:     "invalid: no special characters",
			password: "passWORD12345",
			wantErr:  errInvalidPassNoSpecial,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
