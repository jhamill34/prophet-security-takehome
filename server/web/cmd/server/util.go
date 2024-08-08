package main

import (
	"encoding/json"
	"net/http"
	"net/netip"
	"strconv"
)

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

func ParseIp(val string) netip.Addr {
	if val == "" {
		return netip.IPv4Unspecified()
	}

	addr, err := netip.ParseAddr(val)
	if err != nil {
		panic(err)
	}

	return addr
}

func Json(w http.ResponseWriter, val any, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	return encoder.Encode(val)
}
