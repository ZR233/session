/*
@Time : 2019-07-12 16:39
@Author : zr
*/
package model

import "time"

type Session struct {
	Token      string
	UserId     string
	Channel    string
	ExpireAt   time.Time
	JsonFields interface{}
}
