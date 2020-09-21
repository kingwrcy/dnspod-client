package main

import (
	"github.com/lxn/walk"
)

func ShowErrMsg(msg string) {
	walk.MsgBox(nil, "错误", msg, walk.MsgBoxIconError)
}
