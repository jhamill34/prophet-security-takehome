package main

import (
	"net/http"
	"slices"
	"strings"
	"time"

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
	r.With(PaginatedContext).Get("/", n.ListFilteredNodes())
	return r
}

func (n *NodeResource) ListFilteredNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		pagination := req.Context().Value("pagination").(PaginatedInput)
		cursor := ParseIp(pagination.Cursor)

		allowListIdStr := req.URL.Query().Get("allowlistId")
		invert := req.URL.Query().Get("invert")

		resultMap := make(map[string]*NodeEntry, 0)
		if allowListIdStr == "" {
			dbResult, err := n.queries.ListAllNodes(req.Context(), database.ListAllNodesParams{
				IpAddr: cursor,
				Limit:  pagination.Limit,
			})
			if err != nil {
				InternalServerError(req, resp, err)
				return
			}

			for _, r := range dbResult {
				sourceEntry := NodeSourceEntry{
					SourceId:      r.SourceID,
					Version:       r.Version.Int64,
					LastExecution: r.LastExecution.Time.Format(time.RFC3339),
				}
				if entry, ok := resultMap[r.IpAddr.String()]; ok {
					entry.Sources = append(entry.Sources, sourceEntry)
				} else {
					resultMap[r.IpAddr.String()] = &NodeEntry{
						IpAddr:  r.IpAddr.String(),
						Sources: []NodeSourceEntry{sourceEntry},
					}
				}
			}
		} else {
			allowListId, err := AssertInt(allowListIdStr)
			if err != nil {
				Err(req, resp, "Expected allowListId to be an integer", 400, err)
				return
			}

			if invert != "true" {
				dbResult, err := n.queries.ListNodesWithoutAllowlist(req.Context(), database.ListNodesWithoutAllowlistParams{
					IpAddr: cursor,
					Limit:  pagination.Limit,
					ListID: allowListId,
				})
				if err != nil {
					InternalServerError(req, resp, err)
					return
				}

				for _, r := range dbResult {
					sourceEntry := NodeSourceEntry{
						SourceId:      r.SourceID,
						Version:       r.Version.Int64,
						LastExecution: r.LastExecution.Time.Format(time.RFC3339),
					}
					if entry, ok := resultMap[r.IpAddr.String()]; ok {
						entry.Sources = append(entry.Sources, sourceEntry)
					} else {
						resultMap[r.IpAddr.String()] = &NodeEntry{
							IpAddr:  r.IpAddr.String(),
							Sources: []NodeSourceEntry{sourceEntry},
						}
					}
				}
			} else {
				dbResult, err := n.queries.ListFilteredAllowlistNodes(req.Context(), database.ListFilteredAllowlistNodesParams{
					IpAddr: cursor,
					Limit:  pagination.Limit,
					ListID: allowListId,
				})
				if err != nil {
					InternalServerError(req, resp, err)
					return
				}

				for _, r := range dbResult {
					sourceEntry := NodeSourceEntry{
						SourceId:      r.SourceID,
						Version:       r.Version.Int64,
						LastExecution: r.LastExecution.Time.Format(time.RFC3339),
					}
					if entry, ok := resultMap[r.IpAddr.String()]; ok {
						entry.Sources = append(entry.Sources, sourceEntry)
					} else {
						resultMap[r.IpAddr.String()] = &NodeEntry{
							IpAddr:  r.IpAddr.String(),
							Sources: []NodeSourceEntry{sourceEntry},
						}
					}
				}
			}
		}

		result := make([]NodeEntry, 0, len(resultMap))
		for _, v := range resultMap {
			result = append(result, *v)
		}

		slices.SortFunc(result, func(a, b NodeEntry) int {
			return strings.Compare(a.IpAddr, b.IpAddr)
		})

		paginated := MakePaginated(result, int(pagination.Limit), func(node NodeEntry) string {
			return node.IpAddr
		})

		Json(req, resp, paginated, 200)
	}
}
