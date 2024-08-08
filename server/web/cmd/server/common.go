package main

type PaginatedInput struct {
	Cursor string
	Limit  int32
}

type Paginated[T any] struct {
	Cursor  string `json:"cursor"`
	Total   int    `json:"total"`
	HasMore bool   `json:"has_more"`
	Data    []T    `json:"data"`
}

type NodeEntry struct {
	IpAddr  string            `json:"ip_addr"`
	Sources []NodeSourceEntry `json:"sources"`
}

type NodeSourceEntry struct {
	SourceId      int32  `json:"source_id"`
	Version       int64  `json:"version"`
	LastExecution string `json:"last_execution"`
}

type SourceEntry struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	Period        string `json:"period"`
	LastExecution string `json:"last_execution"`
	Version       int64  `json:"version"`
	Running       bool   `json:"running"`
}

type AllowlistEntry struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type AllowlistEntryItem struct {
	ID     int32  `json:"id"`
	Cidr   string `json:"cidr_block"`
	ListID int32  `json:"allowlist_id"`
}

type ErrorMessage struct {
	Message   string `json:"msg"`
	RequestId string `json:"request_id"`
}
