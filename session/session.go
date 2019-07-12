/*
@Time : 2019-07-12 16:39
@Author : zr
*/
package session

type Session struct {
	Token      string
	UserId     string
	Channel    string
	JsonFields interface{}
}
