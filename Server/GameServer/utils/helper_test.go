package utils

import (
	"testing"
	"time"
)

func TestUnmarshalItemsInt(t *testing.T) {
	ok :=  GetNowTimeStamp() == time.Now().Unix()
	println(ok)

}
