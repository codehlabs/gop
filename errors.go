package gop

import "errors"

var (
	ErrUnabelToAuthenticate = errors.New("unable to authenticate password or username does not match")
)
