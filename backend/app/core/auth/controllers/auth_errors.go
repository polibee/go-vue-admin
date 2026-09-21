package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

var errMissingToken = errors.New("missing authorization token")

func parseAuthToken(ctx http.Context) error {
	token := ctx.Request().Header("Authorization")
	if token == "" {
		return errMissingToken
	}
	_, err := facades.Auth(ctx).Parse(token)
	return err
}
