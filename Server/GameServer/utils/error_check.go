package utils

import (
	"battle_rabbit/service/log"
	"fmt"
)

// print err and return pass false if not nil
func ErrCheck(err error, msg string, param ...interface{}) (out bool) {
	if err != nil {
		log.Error(err,fmt.Sprintf(msg,param...))
		out = true
	}
	return
}


// print err if not nil
func ErrPrint(err error, msg string, param ...interface{}) {
	if err != nil {
		log.Error(err,fmt.Sprintf(msg,param...))
	}
}
