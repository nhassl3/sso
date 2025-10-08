package tests

import (
	"testing"

	"github.com/nhassl3/sso-app/internals/storage"
	"github.com/nhassl3/sso-app/tests/suite"
	ssov1 "github.com/nhassl3/sso-contracts/generated/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoubleRegistration(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email, password := st.NewEmail(), st.NewPassword()

	RespReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, RespReg.GetUserId())

	RespReg2, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, RespReg2.GetUserId())
	assert.ErrorContains(t, err, storage.ErrUserExists.Error())

	RestLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    suite.AppID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, RestLogin.GetToken())
}
