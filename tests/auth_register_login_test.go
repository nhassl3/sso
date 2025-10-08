package tests

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhassl3/sso-app/tests/suite"
	ssov1 "github.com/nhassl3/sso-contracts/generated/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email, password := st.NewEmail(), st.NewPassword()

	RespReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, RespReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    suite.AppID,
	})
	require.NoError(t, err)

	jwtTokenAsserts(t, respLogin, RespReg.GetUserId(), email, st.Cfg.TokenTTL)
}

func jwtTokenAsserts(t *testing.T, resp *ssov1.LoginResponse, UID int64, email string, tokenTTL time.Duration) {
	loginTime := time.Now()

	token := resp.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(suite.AppSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, UID, int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, suite.AppID, int32(claims["app_id"].(float64)))

	assert.InDelta(t, loginTime.Add(tokenTTL).Unix(), claims["exp"].(float64), suite.DeltaSecond)
}
