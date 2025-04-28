package lib

import (
	"fmt"
	"strconv"
)

func AnyToInt(s any) int {
	str := fmt.Sprint(s)
	dataint, _ := strconv.Atoi(str)
	return dataint
}


func AnyToUInt(s any) uint {
	return uint(AnyToInt(s))
}