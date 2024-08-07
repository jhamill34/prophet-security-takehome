package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	r.Get("/all", n.ListAllNodes())
	r.Get("/", n.ListFilteredNodes())
	return r
}

func (n *NodeResource) ListAllNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		cursor := req.URL.Query().Get("after")
		limitStr := req.URL.Query().Get("limit")

		var limit int64
		var err error
		if limitStr == "" {
			limit = 10
		} else {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				panic(err)
			}
		}

		result, err := n.queries.ListAllExistingNodes(req.Context(), database.ListAllExistingNodesParams{
			IpAddr: cursor,
			Limit:  int32(limit),
		})

		if err != nil {
			panic(err)
		}

		marshaled, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(200)
		resp.Write(marshaled)
	}
}

func (n *NodeResource) ListFilteredNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		result, err := n.queries.ListAllExistingNodes(req.Context(), database.ListAllExistingNodesParams{
			IpAddr: "",
			Limit:  10,
		})

		if err != nil {
			panic(err)
		}

		marshaled, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(200)
		resp.Write(marshaled)
	}
}
