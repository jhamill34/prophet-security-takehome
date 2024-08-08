package main

import (
	"encoding/json"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

type NodeResource struct {
	queries *database.Queries
}

func NewNodeResource(queries *database.Queries) *NodeResource {
	return &NodeResource{queries}
}

func (r *NodeResource) Path() string {
	return "/nodes"
}

func (n *NodeResource) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", n.ListFilteredNodes())
	return r
}

func (n *NodeResource) ListFilteredNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var err error

		cursor := ParseIp(req.URL.Query().Get("after"))
		allowListIdStr := req.URL.Query().Get("allowlistId")
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		var result []netip.Addr
		if allowListIdStr == "" {
			result, err = n.queries.ListAllNodes(req.Context(), database.ListAllNodesParams{
				IpAddr: cursor,
				Limit:  limit,
			})
			if err != nil {
				panic(err)
			}
		} else {
			allowListId := AssertInt(allowListIdStr)
			result, err = n.queries.ListFilteredAllowlistNodes(req.Context(), database.ListFilteredAllowlistNodesParams{
				IpAddr: cursor,
				Limit:  limit,
				ListID: allowListId,
			})

			if err != nil {
				panic(err)
			}
		}

		marshaled, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(200)
		resp.Write(marshaled)
	}
}
