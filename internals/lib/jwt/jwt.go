package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhassl3/sso-app/internals/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":  user.Email,
		"exp":    time.Now().Add(duration).Unix(),
		"uid":    user.ID,
		"app_id": app.ID,
	}).SignedString([]byte(app.Secret))
}
