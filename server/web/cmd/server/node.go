package main

import (
	"encoding/json"
	"net/http"

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
	r.Get("/", n.ListNodes())
	return r
}

func (n *NodeResource) ListNodes() http.HandlerFunc {
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
