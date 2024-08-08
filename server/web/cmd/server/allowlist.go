package main

import (
	"encoding/json"
	"net/http"
	"net/netip"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

type AllowListResource struct {
	queries *database.Queries
}

func NewAllowListResource(queries *database.Queries) *AllowListResource {
	return &AllowListResource{
		queries,
	}
}

func (a *AllowListResource) Path() string {
	return "/allowlist"
}

func (a *AllowListResource) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", a.ListAllLists())
	r.Post("/", a.CreateAllowList())
	r.Delete("/{id}", a.DeleteAllowList())

	r.Get("/{id}/entry", a.ListAllowList())
	r.Post("/{id}/entry", a.AddToList())
	r.Delete("/{id}/entry/{entryId}", a.RemoveFromList())

	return r
}

func (a *AllowListResource) ListAllLists() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		after := ParseIntDefault(req.URL.Query().Get("after"), -1)
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		result, err := a.queries.ListAllLists(req.Context(), database.ListAllListsParams{
			ID:    after,
			Limit: limit,
		})

		if err != nil {
			panic(err)
		}

		resultBytes, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(200)
		resp.Write(resultBytes)
	}
}

func (a *AllowListResource) ListAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id := AssertInt(chi.URLParam(req, "id"))

		entries, err := a.queries.ListEntriesForAllowList(req.Context(), id)
		if err != nil {
			panic(err)
		}

		result, err := json.Marshal(entries)
		if err != nil {
			panic(err)
		}
		resp.WriteHeader(200)
		resp.Write(result)
	}
}

type CreateAllowListInput struct {
	Name string `json:"name"`
}

func (a *AllowListResource) CreateAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var input CreateAllowListInput
		json.NewDecoder(req.Body).Decode(&input)
		list, err := a.queries.CreateAllowList(req.Context(), input.Name)
		if err != nil {
			panic(err)
		}

		respBytes, err := json.Marshal(list)
		if err != nil {
			panic(err)
		}
		resp.WriteHeader(201)
		resp.Write(respBytes)
	}
}

func (a *AllowListResource) DeleteAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		idStr := chi.URLParam(req, "id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			panic(err)
		}

		// TODO: Delete all associated entries

		err = a.queries.DeleteAllowList(req.Context(), int32(id))
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}
}

type ListEntryInput struct {
	IpAddr string `json:"ip_addr"`
}

func (a *AllowListResource) AddToList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		listId := AssertInt(chi.URLParam(req, "id"))

		var input ListEntryInput
		json.NewDecoder(req.Body).Decode(&input)

		ipAddr, err := netip.ParsePrefix(input.IpAddr)
		if err != nil {
			panic(err)
		}

		entry, err := a.queries.AddToAllowlist(req.Context(), database.AddToAllowlistParams{
			IpAddr: ipAddr,
			ListID: listId,
		})
		if err != nil {
			panic(err)
		}

		entryBytes, err := json.Marshal(entry)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(201)
		resp.Write(entryBytes)
	}
}

func (a *AllowListResource) RemoveFromList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		listId := AssertInt(chi.URLParam(req, "id"))
		entryId := AssertInt(chi.URLParam(req, "entryId"))

		err := a.queries.RemoveFromAllowlist(req.Context(), database.RemoveFromAllowlistParams{
			ListID: listId,
			ID:     entryId,
		})
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}

}
