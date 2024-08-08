package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/netip"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

func AssertInt(val string) (int32, error) {
	num, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(num), nil
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

func Json(r *http.Request, w http.ResponseWriter, val any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(val)
	if err != nil {
		Err(r, w, "Internal Server Error", 500, err)
	}
}

func InternalServerError(r *http.Request, w http.ResponseWriter, err error) {
	Err(r, w, "Internal Server Error", 500, err)
}

func Err(r *http.Request, w http.ResponseWriter, msg string, status int, err error) {
	requestId := r.Context().Value(middleware.RequestIDKey).(string)
	if err != nil {
		logInternalError(requestId, msg, err, false)
	} else {
		logInternalError(requestId, msg, errors.New("UNKNOWN"), false)
	}

	message := ErrorMessage{
		Message:   msg,
		RequestId: requestId,
	}
	Json(r, w, message, status)
}

func logInternalError(requestId string, msg string, err error, shoudlPanic bool) {
	if shoudlPanic {
		slog.Error(msg, slog.String("request_id", requestId))
		panic(err)
	} else {
		slog.Error(msg,
			slog.String("internal_error", err.Error()),
			slog.String("request_id", requestId),
		)
	}
}
