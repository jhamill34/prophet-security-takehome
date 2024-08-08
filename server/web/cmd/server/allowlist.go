package main

import (
	"encoding/json"
	"net/http"
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
	r.Post("/", a.CreateAllowList())
	r.Get("/{id}", a.ListAllowList())
	r.Delete("/{id}", a.DeleteAllowList())
	r.Post("/{id}/add", a.AddToList())
	r.Post("/{id}/remove", a.RemoveFromList())
	return r
}

func (a *AllowListResource) ListAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		idStr := chi.URLParam(req, "id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			panic(err)
		}

		entries, err := a.queries.ListEntriesForAllowList(req.Context(), int32(id))
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
		listIdStr := chi.URLParam(req, "id")
		listId, err := strconv.ParseInt(listIdStr, 10, 32)
		if err != nil {
			panic(err)
		}

		var input ListEntryInput
		json.NewDecoder(req.Body).Decode(&input)

		entry, err := a.queries.AddToAllowlist(req.Context(), database.AddToAllowlistParams{
			IpAddr: input.IpAddr,
			ListID: int32(listId),
		})

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
		listIdStr := chi.URLParam(req, "id")
		listId, err := strconv.ParseInt(listIdStr, 10, 32)
		if err != nil {
			panic(err)
		}

		var input ListEntryInput
		json.NewDecoder(req.Body).Decode(&input)

		err = a.queries.RemoveFromAllowlist(req.Context(), database.RemoveFromAllowlistParams{
			IpAddr: input.IpAddr,
			ListID: int32(listId),
		})
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}

}
