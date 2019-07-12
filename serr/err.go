/*
@Time : 2019-07-12 17:28
@Author : zr
*/
package serr

const (
	_ = iota
	TokenNotFind
	RedisErr
	JsonErr
)

type Error struct {
	error
	Code int
}

func NewErr(err error, code int) error {
	return Error{err, code}
}
