package controllers

import "errors"

var errMissingToken = errors.New("missing authorization token")
