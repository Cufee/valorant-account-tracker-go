package database

import "errors"

var ErrNotFound = errors.New("value with the given table/key not found")
