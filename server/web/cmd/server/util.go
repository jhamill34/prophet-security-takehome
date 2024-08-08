package main

import "strconv"

func AssertInt(val string) int32 {
	num, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		panic(err)
	}

	return int32(num)
}

func ParseIntDefault(val string, def int32) int32 {
	if val == "" {
		return def
	}

	num, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		panic(err)
	}

	return int32(num)
}
