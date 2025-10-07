package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nhassl3/sso-app/internals/domain/models"
	njwt "github.com/nhassl3/sso-app/internals/lib/jwt"
	"github.com/nhassl3/sso-app/internals/lib/logger/sl"
	"github.com/nhassl3/sso-app/internals/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	opLogin           = "auth.Login"
	opRegisterNewUser = "auth.RegisterNewUser"
	opIsAdmin         = "auth.IsAdmin"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid application ID")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// NewAuth returns a new instance of the Auth service
func NewAuth(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, hashPassword []byte) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (app models.App, err error)
}

// Login checks if user with given credentials exists in the system.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error
func (a *Auth) Login(ctx context.Context, email string, password string, appID int32) (token string, err error) {
	log := a.log.With(slog.String("op", opLogin))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("failed to found user in the system", sl.Err(err))

			return "", sl.ErrUpLevel(opLogin, ErrInvalidCredentials.Error())
		}

		log.Error("failed to get user", sl.Err(err))

		return "", sl.ErrUpLevel(opLogin, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password)); err != nil {
		log.Info(ErrInvalidCredentials.Error(), sl.Err(err))

		return "", sl.ErrUpLevel(opLogin, ErrInvalidCredentials.Error())
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("failed to found app in the system", sl.Err(err))

			return "", sl.ErrUpLevel(opLogin, ErrInvalidAppID.Error())
		}

		log.Error("failed to get app", sl.Err(err))

		return "", sl.ErrUpLevel(opLogin, err.Error())
	}

	token, err = njwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))

		return "", sl.ErrUpLevel(opLogin, err.Error())
	}

	return
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
	log := a.log.With(slog.String("op", opRegisterNewUser))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, sl.ErrUpLevel(opRegisterNewUser, err.Error())
	}

	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return 0, sl.ErrUpLevel(opRegisterNewUser, ErrUserExists.Error())
		}

		log.Error("failed to save user in database", sl.Err(err))

		return 0, sl.ErrUpLevel(opRegisterNewUser, err.Error())
	}

	return
}

// IsAdmin checks if the user has administrator rights on the system.
// If user doesn't have, returns false else true.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	log := a.log.With(slog.String("op", opIsAdmin))

	isAdmin, err = a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("failed to found app in the system", sl.Err(err))

			return false, sl.ErrUpLevel(opIsAdmin, ErrInvalidAppID.Error())
		}

		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("failed to found user in the system", sl.Err(err))

			return false, sl.ErrUpLevel(opIsAdmin, ErrInvalidUserID.Error())
		}

		log.Error("failed to get user", sl.Err(err))

		return false, sl.ErrUpLevel(opIsAdmin, err.Error())
	}

	return
}
