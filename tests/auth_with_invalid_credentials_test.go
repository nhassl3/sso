package tests

import (
	"testing"

	"github.com/nhassl3/sso-app/tests/suite"
	ssov1 "github.com/nhassl3/sso-contracts/generated/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidCredentials_Registration(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		Name        string
		Email       string
		Password    string
		ExceptError string
	}{
		{
			Name:        "With empty Email field",
			Email:       "",
			Password:    st.NewPassword(),
			ExceptError: "value length must be",
		},
		{
			Name:        "With empty Password field",
			Email:       st.NewEmail(),
			Password:    "",
			ExceptError: "value length must be",
		},
		{
			Name:        "With invalid Email and Password fields",
			Email:       "",
			Password:    "",
			ExceptError: "value length must be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.Email,
				Password: tt.Password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.ExceptError)
		})
	}
}

func TestInvalidCredentials_Login(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	tests := []struct {
		Name        string
		Email       string
		Password    string
		AppID       int32
		ExceptError string
	}{
		{
			Name:        "With empty Email field",
			Email:       "",
			Password:    st.NewPassword(),
			AppID:       suite.AppID,
			ExceptError: "value length must be",
		},
		{
			Name:        "With empty Password field",
			Email:       st.NewEmail(),
			Password:    "",
			AppID:       suite.AppID,
			ExceptError: "value length must be",
		},
		{
			Name:        "With invalid Email and Password fields",
			Email:       "",
			Password:    "",
			AppID:       suite.AppID,
			ExceptError: "value length must be",
		},
		{
			Name:        "With invalid App ID field",
			Email:       st.NewEmail(),
			Password:    st.NewPassword(),
			AppID:       suite.EmptyAppID,
			ExceptError: "AppID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    st.NewEmail(),
				Password: st.NewPassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.Email,
				Password: tt.Password,
				AppId:    tt.AppID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.ExceptError)
		})
	}
}

func TestValidCredentials_Registration(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	t.Run("Valid registration credentials", func(t *testing.T) {
		respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    st.NewEmail(),
			Password: st.NewPassword(),
		})
		require.NoError(t, err)
		assert.NotEmpty(t, respReg.GetUserId())
	})
}

func TestValidCredentials_Login(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email, password := st.NewEmail(), st.NewPassword()

	t.Run("Valid registration credentials", func(t *testing.T) {
		respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, respReg.GetUserId())

		respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Email:    email,
			Password: password,
			AppId:    suite.AppID,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, respLogin.GetToken())

		jwtTokenAsserts(t, respLogin, respReg.GetUserId(), email, st.Cfg.TokenTTL)
	})
}
