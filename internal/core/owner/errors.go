package owner

import (
	"errors"
)

var (
	ErrNotOwner = errors.New("apenas o dono pode efetuar essa ação")
)
