package routes

import (
	"net/netip"

	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
)

type CursorExtractor[T any] func(item T) string

func MakePaginated[T any](data []T, limit int, cursorExt CursorExtractor[T]) api.PaginatedMetadata {
	total := len(data)
	hasMore := total >= limit

	cursor := ""
	if total > 0 {
		cursor = cursorExt(data[total-1])
	}

	return api.PaginatedMetadata{
		Cursor:  cursor,
		Total:   total,
		HasMore: hasMore,
	}
}

func DefaultValue[T any](ptr *T, value T) T {
	if ptr == nil {
		return value
	}

	return *ptr
}

func ParseIp(val *string) (netip.Addr, error) {
	if val == nil {
		return netip.IPv4Unspecified(), nil
	}

	return netip.ParseAddr(*val)
}
